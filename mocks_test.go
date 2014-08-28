package batcher

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
)

type MockFlusherImplementation struct {
	flushed chan Flushable
	delay   time.Duration
}

func NewMockFlusherImplementation() *MockFlusherImplementation {
	mf := new(MockFlusherImplementation)
	mf.flushed = make(chan Flushable)
	mf.delay = time.Second * 0
	return mf
}

func (this *MockFlusherImplementation) Flush(to_flush Flushable) error {
	time.Sleep(this.delay)
	go func() {
		this.flushed <- to_flush
	}()
	return nil
}

type MockFlusher struct {
	flushed chan Flushable
}

func NewMockFlusher() *MockFlusher {
	mf := new(MockFlusher)
	mf.flushed = make(chan Flushable)
	return mf
}

func (this *MockFlusher) Shutdown() {
	// panic("shutdown not implemented on the mock")
}
func (this *MockFlusher) Flush(to_flush Flushable) error {
	go func() {
		this.flushed <- to_flush
	}()
	return nil
}

func ShouldReceiveStringIn(c chan string, r string, d time.Duration) {
	select {
	case a := <-c:
		So(a, ShouldEqual, r)
	case <-time.After(d):
		So(false, ShouldBeTrue)
	}
}

func ShouldReceiveIn(c chan interface{}, r interface{}, d time.Duration) {
	select {
	case a := <-c:
		So(a, ShouldResemble, r)
	case <-time.After(d):
		So(false, ShouldBeTrue)
	}
}

func ShouldReceiveSliceIn(c chan []interface{}, r []interface{}, d time.Duration) {
	select {
	case a := <-c:
		So(a, ShouldResemble, r)
	case <-time.After(d):
		So(false, ShouldBeTrue)
	}
}

func ShouldNotReceiveSliceIn(c chan []interface{}, d time.Duration) {
	select {
	case <-c:
		So(false, ShouldBeTrue)
	case <-time.After(d):
	}
}

func ShouldReceiveFlushableStringsIn(c chan Flushable, r []string, d time.Duration) {
	select {
	case a := <-c:
		So(a.Strings(), ShouldResemble, r)
	case <-time.After(d):
		So(false, ShouldBeTrue)
	}
}

func ShouldReceiveFlushableIn(c chan Flushable, r Flushable, d time.Duration) {
	select {
	case a := <-c:
		So(a, ShouldResemble, r)
	case <-time.After(d):
		So(false, ShouldBeTrue)
	}
}

func ShouldNotReceiveFlushableIn(c chan Flushable, d time.Duration) {
	select {
	case <-c:
		So(false, ShouldBeTrue)
	case <-time.After(d):
	}
}

type MockRedisPool struct {
	dos           chan string
	CommandsDelay time.Duration
}

func NewMockRedisPool() *MockRedisPool {
	m := new(MockRedisPool)
	m.dos = make(chan string, 10)
	return m
}

func (this *MockRedisPool) Get() redis.Conn {
	rc := NewMockRedisConnection(this.dos)
	if int64(this.CommandsDelay) != 0 {
		rc.CommandsDelay = this.CommandsDelay
	}
	return rc
}

type MockRedisConnection struct {
	dos           chan string
	CommandsDelay time.Duration
}

func NewMockRedisConnection(dos chan string) *MockRedisConnection {
	m := new(MockRedisConnection)
	m.dos = dos
	return m
}

func (this *MockRedisConnection) Close() error {
	return nil
}

func (this *MockRedisConnection) Err() error {
	return nil
}

func (this *MockRedisConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if int64(this.CommandsDelay) != 0 {
		time.Sleep(this.CommandsDelay)
	}
	s := make([]interface{}, len(args)+1)
	s[0] = commandName
	for i, v := range args {
		s[i+1] = v
	}
	encoded, _ := json.Marshal(s)
	this.dos <- string(encoded)
	return nil, nil
}

func (this *MockRedisConnection) Send(commandName string, args ...interface{}) error {
	return nil
}

func (this *MockRedisConnection) Flush() error {
	return nil
}

func (this *MockRedisConnection) Receive() (reply interface{}, err error) {
	return nil, nil
}

type MockRedisPoolDiscard struct {
	dos chan string
}

func NewMockRedisPoolDiscard() *MockRedisPoolDiscard {
	m := new(MockRedisPoolDiscard)
	m.dos = make(chan string, 10)
	return m
}

func (this *MockRedisPoolDiscard) Get() redis.Conn {
	rc := NewMockRedisConnectionDiscard(this.dos)
	return rc
}

type MockRedisConnectionDiscard struct {
	dos chan string
}

func NewMockRedisConnectionDiscard(dos chan string) *MockRedisConnectionDiscard {
	m := new(MockRedisConnectionDiscard)
	m.dos = dos
	return m
}

func (this *MockRedisConnectionDiscard) Close() error {
	return nil
}

func (this *MockRedisConnectionDiscard) Err() error {
	return nil
}

func (this *MockRedisConnectionDiscard) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return nil, nil
}

func (this *MockRedisConnectionDiscard) Send(commandName string, args ...interface{}) error {
	return nil
}

func (this *MockRedisConnectionDiscard) Flush() error {
	return nil
}

func (this *MockRedisConnectionDiscard) Receive() (reply interface{}, err error) {
	return nil, nil
}

type MockStruct struct {
	Name string
}

func (this MockStruct) String() string {
	encoded, _ := json.Marshal(this)
	return string(encoded)
}
