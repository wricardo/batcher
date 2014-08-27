package function_flusher

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
	"github.com/wricardo/batcher"
	"encoding/json"
)
func init(){
	fmt.Sprint("")
}

func TestIfFunctionIsCalledWhenFlusherFlushes (t *testing.T) {
	Convey("Given a function flusher", t, func() {
		flushed := make(chan []string,0)
		f := func(to_flush batcher.Flushable) error{
			flushed <- to_flush.Strings()
			return nil
		}

		ff := NewFunctionFlusher(0, f)
		flusher := batcher.NewDefaultFlusher(3, ff)

		Convey("When Flushing one flushable struct", func() {
			to_flush := make(batcher.Flushable, 1)
			to_flush[0] = TestStruct{"name1", 1987}
			err := flusher.Flush(to_flush)

			Convey("No error should have occurred", func() {
				So(err, ShouldBeNil)
			})
			Convey("Flusher should call Flush on FunctionFlusher", func() {
				expected := make([]string,1)
				expected[0] = "{\"Field1\":\"name1\",\"Field2\":1987}"
				So(<-flushed, ShouldResemble, expected)
			})
		})

		Convey("When flushing a flushable slice of structs", func() {
			to_flush := make(batcher.Flushable, 2)
			to_flush[0] = TestStruct{"name1", 1987}
			to_flush[1] = TestStruct{"name2", 1990}

			err := flusher.Flush(to_flush)

			Convey("No error should have occurred", func() {
				So(err, ShouldBeNil)
			})
			Convey("The flushed values should match", func() {
				expected := make([]string,2)
				expected[0] = "{\"Field1\":\"name1\",\"Field2\":1987}"
				expected[1] = "{\"Field1\":\"name2\",\"Field2\":1990}"
				So(<-flushed, ShouldResemble, expected)
			})
		})
	})
}

type TestStructs []TestStruct

func (this TestStructs) Strings() []string {
	tmp := make([]string, 0)
	for _, v := range this {
		tmp2 := v.String()
		tmp = append(tmp, tmp2)
	}
	return tmp
}

type TestStruct struct {
	Field1 string
	Field2 int
}

func (this TestStruct) String() string {
	encoded, _ := json.Marshal(this)
	return string(encoded)
}
