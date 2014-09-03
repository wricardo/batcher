package batcher

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringTypeAndMethod(t *testing.T) {
	Convey("Given a String", t, func() {
		s := String("test_string")
		So(s.String(), ShouldEqual,"test_string")
	})
}

func TestIntTypeAndMethod(t *testing.T) {
	Convey("Given a Int", t, func() {
		s := Int(123)
		So(s.String(), ShouldEqual,"123")
	})
}

