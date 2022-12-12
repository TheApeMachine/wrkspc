package errnie

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type ErrorTestContext struct {
	instance error
	msg      string
}

func NewErrorTestContext(err error, msg string) *ErrorTestContext {
	return &ErrorTestContext{
		instance: NewError(err),
		msg:      msg,
	}
}

func TestNewError(t *testing.T) {
	Convey("Given a nil error", t, func() {
		ctx := NewErrorTestContext(nil, "")

		Convey("It should return nil", func() {
			So(ctx.instance, ShouldBeNil)
		})
	})

	Convey("Given an error type", t, func() {
		msg := "this is a test"
		ctx := NewErrorTestContext(
			fmt.Errorf("VALIDATION ERROR|%s", msg), msg,
		)

		Convey("It should return the error message", func() {
			So(ctx.instance.Error(), ShouldEqual, ctx.msg)
		})
	})
}

func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewError(nil)
	}
}
