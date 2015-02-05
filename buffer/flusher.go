package buffer

import (
	"errors"
	"time"
	//	"fmt"
)

type Flusher interface {
	Flush(Flushable) error
	Shutdown()
}

type FlusherImplementation interface {
	Flush(Flushable) error
}

type Flushable []Collectable

func (this Flushable) Strings() []string {
	tmp := make([]string, len(this))
	for k, v := range this {
		tmp[k] = v.String()
	}
	return tmp
}

type DefaultFlusher struct {
	shutdown       chan chan bool
	flush_chan     chan Flushable
	implementation FlusherImplementation
}

func NewDefaultFlusher(flush_buffer_size int, implementation FlusherImplementation) *DefaultFlusher {
	rf := new(DefaultFlusher)
	rf.flush_chan = make(chan Flushable, flush_buffer_size)
	rf.shutdown = make(chan chan bool, 0)
	rf.implementation = implementation
	go rf.listenAndFlush()
	return rf
}

func (this *DefaultFlusher) Shutdown() {
	finished := make(chan bool, 0)
	this.shutdown <- finished
	<-finished
}

func (this *DefaultFlusher) listenAndFlush() {
	var shutdown chan bool
	shutdown_func := func() {
		close(this.flush_chan)
		shutdown <- true
	}
L:
	for {
		select {
		case fs := <-this.shutdown:
			shutdown = fs
			if len(this.flush_chan) == 0 {
				shutdown_func()
				break L
			}
		case to_flush, received := <-this.flush_chan:
			if received == false {
				break L
			}
			this.implementation.Flush(to_flush)
			if len(this.flush_chan) == 0 && shutdown != nil {
				shutdown_func()
				break L
			}
		case <-time.After(time.Second * 1):
			if shutdown != nil {
				shutdown_func()
				break L
			}
		}
	}
}

func (this *DefaultFlusher) Flush(to_flush Flushable) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("This flusher is shut down")
		}
	}()

	select {
	case this.flush_chan <- to_flush:
	case <-time.After(time.Second * 10):
		err = errors.New("Timeout on send")
	}

	return err
}
