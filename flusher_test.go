package batcher

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
//	"encoding/json"
	"runtime"
	"github.com/wricardo/batcher/mock"
)

func init() {
	fmt.Sprint("")
}

func TestShuttingDownFlusherWaitsForAllTheCommandsToBeFinished(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	SkipConvey("Given a flusher with buffer 3", t, func() {
		mfi := NewMockFlusherImplementation()
		mfi.delay = time.Second * 1
		flusher := NewDefaultFlusher(3, mfi)

		Convey("When Flushing 3 structs inline and 2 on goroutines", func() {
			mock_commands_received := make(chan []Flushable, 0)
			go func() {
				x := make([]Flushable, 0)
				for i := 0; i < 5; i++ {
					x = append(x, <-mfi.flushed)
				}
				mock_commands_received <- x
			}()
			tmp := make([]Flushable, 5)
			tmp[0] = make(Flushable, 1)
			tmp[0][0] = mock.TestStruct{"name1", 1}
			tmp[1] = make(Flushable, 1)
			tmp[1][0] = mock.TestStruct{"name2", 2}
			tmp[2] = make(Flushable, 1)
			tmp[2][0] = mock.TestStruct{"name3", 3}
			tmp[3] = make(Flushable, 1)
			tmp[3][0] = mock.TestStruct{"name4", 4}
			tmp[4] = make(Flushable, 1)
			tmp[4][0] = mock.TestStruct{"name5", 5}

			flusher.Flush(tmp[0])
			flusher.Flush(tmp[1])
			flusher.Flush(tmp[2])
			flusher.Flush(tmp[3])
			flusher.Flush(tmp[4])

			Convey("After calling shutting down all 5 messages should have been flushed", func() {
				flusher.Shutdown()
				So(<-mock_commands_received, ShouldResemble, tmp)
			})

		})
	})
}

func TestFlushing(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Convey("Given a flusher", t, func() {
		mfi := NewMockFlusherImplementation()
		flusher := NewDefaultFlusher(3, mfi)

		Convey("When Flushing one flushable struct", func() {
			to_flush := make(Flushable, 1)
			to_flush[0] = mock.TestStruct{"name1", 1987}
			err := flusher.Flush(to_flush)

			Convey("Flusher needs to send a flush command with the flushable struct", func() {
				ShouldReceiveFlushableIn(mfi.flushed, to_flush, time.Duration(time.Second))
			})
			Convey("No error should have occurred", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When flushing a flushable slice of structs", func() {
			to_flush := make(Flushable, 2)
			to_flush[0] = mock.TestStruct{"name1", 1987}
			to_flush[1] = mock.TestStruct{"name2", 1990}

			err := flusher.Flush(to_flush)
			Convey("Flusher needs to call the implementation with the same slice", func() {
				ShouldReceiveFlushableIn(mfi.flushed, to_flush, time.Duration(time.Second))
			})
			Convey("No error should have occurred", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestFlushingAStructMultipleTimes(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Convey("Given a flusher", t, func() {
		mfi := NewMockFlusherImplementation()
		flusher := NewDefaultFlusher(3, mfi)

		Convey("When Flushing a slice with one struct multiple times", func() {
			to_flush := make(Flushable, 1)
			to_flush[0] = mock.TestStruct{"name1", 1987}
			for i := 0; i < 31; i++ {
				go flusher.Flush(to_flush)
			}

			Convey("Flusher needs to call the implementation with the same slice", func() {
				for i := 0; i < 31; i++ {
					ShouldReceiveFlushableIn(mfi.flushed, to_flush, time.Duration(time.Second))
				}
			})
		})

	})
}
