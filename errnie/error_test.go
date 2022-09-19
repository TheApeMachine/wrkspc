package errnie

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewError(errors.New("test error"))
	}
}

func TestNewError(t *testing.T) {
	Convey("Given an error", t, func() {
		err := NewError(errors.New("test error"))

		Convey("It should have the correct type", func() {
			So(err.Type, ShouldEqual, NOK)
		})

		Convey("It should wrap teh error message", func() {
			So(err.Msg, ShouldEqual, "test error")
		})
	})
}
