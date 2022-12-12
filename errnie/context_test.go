package errnie

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type ContextContext struct {
	instance *Context
}

func NewContextContext() *ContextContext {
	return &ContextContext{
		instance: New(),
	}
}

func TestNew(t *testing.T) {
	Convey("Given a new instance", t, func() {
		NewContextContext()
	})
}

func TestTracing(t *testing.T) {
	Convey("Given a new instance", t, func() {
		NewContextContext()

		Convey("And tracing is turned off", func() {
			Tracing(false)

			Convey("It should have the correct setting", func() {
				So(ctx.tracing, ShouldBeFalse)
			})
		})

		Convey("And tracing is turned on", func() {
			Tracing(true)

			Convey("It should have the correct setting", func() {
				So(ctx.tracing, ShouldBeTrue)
			})
		})
	})
}

func TestDebugging(t *testing.T) {
	Convey("Given a new instance", t, func() {
		NewContextContext()

		Convey("And debugging is turned off", func() {
			Debugging(false)

			Convey("It should have the correct setting", func() {
				So(ctx.debugging, ShouldBeFalse)
			})
		})

		Convey("And debugging is turned on", func() {
			Debugging(true)

			Convey("It should have the correct setting", func() {
				So(ctx.debugging, ShouldBeTrue)
			})
		})
	})
}

func BenchmarkNew(b *testing.B) {
	NewContextContext()

	for i := 0; i < b.N; i++ {
		_ = New()
	}
}
