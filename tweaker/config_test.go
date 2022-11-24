package tweaker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

type ConfigContext struct{}

func NewConfigContext() *ConfigContext {
	v := viper.New()
	v.AddConfigPath("../cmd/cfg")
	v.SetConfigType("yml")
	v.SetConfigName(".wrkspc")
	errnie.Handles(v.ReadInConfig())

	return &ConfigContext{}
}

func TestNew(t *testing.T) {
	NewConfigContext()

	Convey("Given a new instance", t, func() {
		cfg := New()

		Convey("It should have an instance of viper", func() {
			So(cfg.v, ShouldNotBeNil)
		})
	})
}

func BenchmarkNew(b *testing.B) {
	NewConfigContext()

	for i := 0; i < b.N; i++ {
		_ = New()
	}
}
