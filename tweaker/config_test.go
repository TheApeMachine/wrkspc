package tweaker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/test"
)

type ConfigContext struct{}

func TestNew(t *testing.T) {
	test.NewConfigContext()

	Convey("Given a new instance", t, func() {
		Convey("It should have an instance of viper", func() {
			So(cfg.v, ShouldNotBeNil)
		})
	})
}

func BenchmarkNew(b *testing.B) {
	test.NewConfigContext()

	for i := 0; i < b.N; i++ {
		_ = New()
	}
}
