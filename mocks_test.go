package batcher

import (
	"encoding/json"
	"time"

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


type MockStruct struct {
	Name string
}

func (this MockStruct) String() string {
	encoded, _ := json.Marshal(this)
	return string(encoded)
}
