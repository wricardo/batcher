package redis_flusher

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
//	"time"
//	"encoding/json"
	"runtime"
	"github.com/wricardo/batcher/mock"
	"github.com/wricardo/batcher/buffer"
)

func init() {
	fmt.Sprint("")
}


func TestLpushFlusherSendsTheCorrectCommandToRedis(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Convey("Given a Mock Redis Pool", t, func() {
		pool := mock.NewMockRedisPool()

		Convey("Given LPushFlusher", func() {
			rpf := NewLPushFlusher("fake_list", 10, pool)

			Convey("Flushing one struct", func() {
				to_flush := make(buffer.Flushable, 1)
				to_flush[0] = mock.TestStruct{"name1", 1987}
				err := rpf.Flush(to_flush)
				So(err, ShouldBeNil)
				tmp := <-pool.Dos
				So(tmp, ShouldEqual,"[\"LPUSH\",\"fake_list\",\"{\\\"Field1\\\":\\\"name1\\\",\\\"Field2\\\":1987}\"]")
			})

			Convey("Flushing two structs", func() {
				to_flush := make(buffer.Flushable, 2)
				to_flush[0] = mock.TestStruct{"name1", 1987}
				to_flush[1] = mock.TestStruct{"name2", 1988}
				err := rpf.Flush(to_flush)
				So(err, ShouldBeNil)
				tmp := <-pool.Dos
				So(tmp, ShouldEqual,"[\"LPUSH\",\"fake_list\",\"{\\\"Field1\\\":\\\"name1\\\",\\\"Field2\\\":1987}\",\"{\\\"Field1\\\":\\\"name2\\\",\\\"Field2\\\":1988}\"]")
			})
		})
	})
}
