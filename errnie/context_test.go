package errnie

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func TestNew(t *testing.T) {
	Convey("Given a new instance", t, func() {
		ctx := New()

		Convey("It should be an AmbientContext", func() {
			So(
				ctx,
				ShouldHaveSameTypeAs,
				AmbientContext{},
			)
		})
	})
}
