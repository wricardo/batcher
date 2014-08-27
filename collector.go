package batcher

import (
	"errors"
	"strconv"
	"time"
)

func init() {
}

type String string

func (this String) String() string {
	return string(this)
}

type Int int

func (this Int) String() string {
	return strconv.Itoa(int(this))
}

type Collectable interface {
	String() string
}

type CollectedMessages []Collectable

func (this CollectedMessages) Strings() []string {
	tmp := make([]string, 0)
	for _, v := range this {
		tmp = append(tmp, v.String())
	}
	return tmp
}

type Collector interface {
	Collect(Collectable) error
}

type bufferedCollector struct {
	Flusher      Flusher
	chan_receive chan Collectable
	shutdown     chan chan bool
	buffer       CollectedMessages
	i            int
	bcc          *bufferedCollectorConfiguration
}

func NewBufferedCollector(bcc *bufferedCollectorConfiguration, flusher Flusher) *bufferedCollector {
	bc := new(bufferedCollector)
	bc.Flusher = flusher
	bc.chan_receive = make(chan Collectable, bcc.bufferSize)
	bc.shutdown = make(chan chan bool, 0)
	bc.buffer = make(CollectedMessages, bcc.bufferSize)
	bc.i = 0
	bc.bcc = bcc
	go bc.receive()
	return bc
}

func (this *bufferedCollector) Shutdown() {
	finished := make(chan bool, 0)
	this.shutdown <- finished
	<-finished
}

func (this *bufferedCollector) Flush() {
	if this.i > 0 {
		to_flush := make(Flushable, this.i)
		copy(to_flush, this.buffer[:this.i])
		this.i = 0
		this.Flusher.Flush(to_flush)
	}
}

func (this *bufferedCollector) receive() {
	var ever_received bool = false
	var shutdown chan bool
	shutdown_func := func() {
		close(this.chan_receive)
		this.Flusher.Shutdown()
		shutdown <- true
	}

L:
	for {
		select {
		case fs := <-this.shutdown:
			shutdown = fs
			if len(this.chan_receive) == 0 {
				if this.i > 0 {
					this.Flush()
				}
				shutdown_func()
				break L
			}
		case r := <-this.chan_receive:
			ever_received = true
			this.CollectUnsafe(r)
			if len(this.chan_receive) == 0 && shutdown != nil {
				this.Flush()
				shutdown_func()
				break L
			}
		case <-time.After(this.bcc.bufferTimeout):
			if ever_received == true {
				this.Flush()
			}
			if shutdown != nil {
				shutdown_func()
				break L
			}
		}
	}
}

func (this *bufferedCollector) CollectUnsafe(r Collectable) {
	this.buffer[this.i] = r
	this.i = this.i + 1
	if this.i >= this.bcc.bufferSize {
		this.Flush()
	}
}

func (this *bufferedCollector) Collect(i Collectable) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("This flusher is shut down")
		}
	}()
	if this.Flusher == nil {
		err = errors.New("Invalid Flusher")
		return
	}
	select {
	case this.chan_receive <- i:
	case <-time.After(time.Second):
		err = errors.New("Timeout on collect")
	}

	return
}

type bufferedCollectorConfiguration struct {
	bufferSize      int
	bufferTimeout   time.Duration
	numberOfWorkers int
}

func NewBufferedCollectorConfiguration(bufferSize int, bufferTimeout time.Duration, numberOfWorkers int) (*bufferedCollectorConfiguration, error) {
	if numberOfWorkers <= 0 {
		return nil, errors.New("Invalid numberOfWorkers")
	}
	if bufferSize <= 0 {
		return nil, errors.New("Invalid bufferSize")
	}
	if int64(bufferTimeout) <= 0 {
		return nil, errors.New("Invalid bufferTimeout")
	}
	bcc := new(bufferedCollectorConfiguration)
	bcc.bufferSize = bufferSize
	bcc.bufferTimeout = bufferTimeout
	bcc.numberOfWorkers = numberOfWorkers
	return bcc, nil
}
