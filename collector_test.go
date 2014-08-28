package batcher

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnsafeSizeBufferedCollector(t *testing.T) {
	Convey("Given a collector", t, func() {
		mf := NewMockFlusher()
		bcc, _ := NewBufferedCollectorConfiguration(2, time.Duration(time.Millisecond*100), 2)
		c := NewBufferedCollector(bcc, mf)

		Convey("When no information is collected", func() {
			Convey("The flusher should not have received any information", func() {
				ShouldNotReceiveFlushableIn(mf.flushed, time.Millisecond*200)
			})
		})

		Convey("When collect once", func() {
			err := c.Collect(String("a"))

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The flusher should NOT have been called within 50 milliseconds", func() {
				ShouldNotReceiveFlushableIn(mf.flushed, time.Millisecond*50)
			})

			Convey("The flusher should have been called within 300 milliseconds", func() {
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"a"}, time.Millisecond*300)
			})
		})

		Convey("When collect until the buffer size is reached", func() {
			c.Collect(String("b"))
			c.Collect(String("c"))

			ShouldReceiveFlushableStringsIn(mf.flushed, []string{"b", "c"}, time.Millisecond*50)
		})

		Convey("When collect until the buffer size is reached with a sleep in the middle", func() {
			c.Collect(String("b"))
			time.Sleep(time.Duration(time.Millisecond * 150))
			c.Collect(String("c"))

			Convey("The flusher should have been called with the first element collected and the flusher should have beem called with the second element collected after the buffer timeout", func() {
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"b"}, time.Millisecond*5)
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"c"}, time.Millisecond*150)
			})
		})

		Convey("When collect twice the size of the buffer", func() {
			c.Collect(String("b"))
			c.Collect(String("c"))
			c.Collect(String("d"))
			c.Collect(String("e"))

			ShouldReceiveFlushableStringsIn(mf.flushed, []string{"b", "c"}, time.Millisecond*50)
			ShouldReceiveFlushableStringsIn(mf.flushed, []string{"d", "e"}, time.Millisecond*50)
		})
	})
}

func TestCollectingAStruct(t *testing.T) {
	Convey("Given a collector", t, func() {
		mf := NewMockFlusher()
		bcc, _ := NewBufferedCollectorConfiguration(2, time.Duration(time.Millisecond*10), 2)
		c := NewBufferedCollector(bcc, mf)

		Convey("When collect a struct", func() {
			c.Collect(MockStruct{Name: "My Name"})

			Convey("The flusher should have been called within 30 milliseconds", func() {
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"{\"Name\":\"My Name\"}"}, time.Millisecond*30)
			})
		})

	})
}

func TestSimpleSizeBufferedCollector(t *testing.T) {
	Convey("Given a collector", t, func() {
		mf := NewMockFlusher()
		bcc, _ := NewBufferedCollectorConfiguration(2, time.Duration(time.Millisecond*100), 2)
		c := NewBufferedCollector(bcc, mf)

		Convey("When no information is collected", func() {
			Convey("The flusher should not have received any information", func() {
				ShouldNotReceiveFlushableIn(mf.flushed, time.Millisecond*200)
			})
		})

		Convey("When collect once", func() {
			err := c.Collect(String("a"))

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The flusher should NOT have been called within 50 milliseconds", func() {
				ShouldNotReceiveFlushableIn(mf.flushed, time.Millisecond*50)
			})

			Convey("The flusher should have been called within 300 milliseconds", func() {
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"a"}, time.Millisecond*300)
			})
		})

		Convey("When collect until the buffer size is reached", func() {
			c.Collect(String("b"))
			c.Collect(String("c"))

			ShouldReceiveFlushableStringsIn(mf.flushed, []string{"b", "c"}, time.Millisecond*50)
		})

		Convey("When collect until the buffer size is reached with a sleep in the middle", func() {
			c.Collect(String("b"))
			time.Sleep(time.Duration(time.Millisecond * 150))
			c.Collect(String("c"))

			Convey("The flusher should have been called with the first element collected and the flusher should have beem called with the second element collected after the buffer timeout", func() {
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"b"}, time.Millisecond*5)
				ShouldReceiveFlushableStringsIn(mf.flushed, []string{"c"}, time.Millisecond*150)
			})
		})

		Convey("When collect twice the size of the buffer", func() {
			c.Collect(String("b"))
			c.Collect(String("c"))
			c.Collect(String("d"))
			c.Collect(String("e"))

			ShouldReceiveFlushableStringsIn(mf.flushed, []string{"b", "c"}, time.Millisecond*50)
			ShouldReceiveFlushableStringsIn(mf.flushed, []string{"d", "e"}, time.Millisecond*50)
		})
	})
}

func TestWaitForAllCollectedMessagesToBeFlushed(t *testing.T) {
	Convey("Given a collector", t, func() {
		mf := NewMockFlusher()
		bcc, _ := NewBufferedCollectorConfiguration(2, time.Duration(time.Millisecond*100), 2)
		c := NewBufferedCollector(bcc, mf)

		Convey("When collect once", func() {
			mock_commands_received := make(chan string, 0)
			go func() {
				var x Flushable
				x = <-mf.flushed
				mock_commands_received <- fmt.Sprint(x)
			}()
			c.Collect(String("a"))

			Convey("AND ask for the collector to shutdown", func() {
				c.Shutdown()

				Convey("The collected message should have been flushed, even if the buffer size hasn't been reached", func() {
					So(<-mock_commands_received, ShouldEqual, "[a]")
				})
			})
		})
	})
}

func TestWaitForAllCollectedMessagesToBeFlushed2(t *testing.T) {
	Convey("Given a collector", t, func() {
		mf := NewMockFlusher()
		bcc, _ := NewBufferedCollectorConfiguration(2, time.Duration(time.Millisecond*100), 2)
		c := NewBufferedCollector(bcc, mf)

		Convey("When collect once", func() {
			mock_commands_received := make(chan string, 2)
			go func() {
				var x Flushable
				x = <-mf.flushed
				mock_commands_received <- fmt.Sprint(x)
				x = <-mf.flushed
				mock_commands_received <- fmt.Sprint(x)
			}()
			c.Collect(String("a"))
			c.Collect(String("b"))
			c.Collect(String("c"))

			Convey("AND ask for the collector to shutdown", func() {
				c.Shutdown()

				Convey("The collected message should have been flushed", func() {
					So(<-mock_commands_received, ShouldEqual, "[a b]")
					So(<-mock_commands_received, ShouldEqual, "[c]")
				})
			})
		})
	})
}

func TestStressCollector(t *testing.T) {
	start := 2
	incr := 100
	for i := start; i < 1000; i = i + incr {
		Convey("Given a valid collector with buffer_size="+strconv.Itoa(i), t, func() {
			mf := NewMockFlusher()
			bcc, _ := NewBufferedCollectorConfiguration(i, time.Duration(time.Second), 2)
			c := NewBufferedCollector(bcc, mf)

			Convey("When no information is collected", func() {
				Convey("The flusher should not have received any information", func() {
					ShouldNotReceiveFlushableIn(mf.flushed, time.Millisecond*200)
				})
			})

			Convey("When collect once", func() {
				c.Collect(String("a"))

				Convey("The flusher should not have been called", func() {
					ShouldNotReceiveFlushableIn(mf.flushed, time.Millisecond*100)
				})
			})

			Convey("When collect until the buffer size is reached", func() {
				expected := make([]string, i)
				for x := 0; x < i; x++ {
					expected[x] = "" + strconv.Itoa(x) + ""
					c.Collect(String(strconv.Itoa(x)))
				}
				ShouldReceiveFlushableStringsIn(mf.flushed, expected, time.Millisecond*100)
			})

			Convey("When collect twice the size of the buffer", func() {
				expected := make([]string, i*2)
				for x := 0; x < i*2; x++ {
					expected[x] = "" + strconv.Itoa(x) + ""
					c.Collect(String(strconv.Itoa(x)))
				}
				ShouldReceiveFlushableStringsIn(mf.flushed, expected[0:i], time.Millisecond*100)
				ShouldReceiveFlushableStringsIn(mf.flushed, expected[i:i*2], time.Millisecond*100)
			})
		})
	}
}

func TestBufferCollectorCreation(t *testing.T) {
	Convey("When constructing a bufferedCollectorConfiguration", t, func() {
		var flagtests = []struct {
			test_case       string
			bufferSize      int
			bufferTimeout   time.Duration
			numberOfWorkers int
			assertBcc       func(interface{}, ...interface{}) string
			assertError     func(interface{}, ...interface{}) string
		}{
			{"Valid parameters", 1, time.Duration(time.Second), 1, ShouldNotBeNil, ShouldBeNil},
			{"Invalid bufferSize 0", 0, time.Duration(time.Second), 1, ShouldBeNil, ShouldNotBeNil},
			{"Invalid bufferSize -1", -1, time.Duration(time.Second), 1, ShouldBeNil, ShouldNotBeNil},
			{"Invalid number of workers", 1, time.Duration(time.Second), 0, ShouldBeNil, ShouldNotBeNil},
			{"Invalid number of workers", 1, time.Duration(time.Second), -1, ShouldBeNil, ShouldNotBeNil},
			{"Invalid bufferTimeout", 1, time.Duration(time.Second * 0), 1, ShouldBeNil, ShouldNotBeNil},
			{"Invalid bufferTimeout", 1, time.Duration(time.Second * -1), 1, ShouldBeNil, ShouldNotBeNil},
		}

		for _, s := range flagtests {
			Convey(s.test_case, func() {
				bufferSize := s.bufferSize
				bufferTimeout := s.bufferTimeout
				numberOfWorkers := s.numberOfWorkers

				Convey("Assertion", func() {
					bcc, err := NewBufferedCollectorConfiguration(bufferSize, bufferTimeout, numberOfWorkers)
					So(bcc, s.assertBcc)
					So(err, s.assertError)
				})
			})
		}
	})
}
