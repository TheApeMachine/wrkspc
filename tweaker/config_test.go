package tweaker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

type ConfigContext struct{}

func NewConfigContext() *ConfigContext {
	// Instantiate viper manually since we are not going through the
	// CLI while running tests.
	viper.AddConfigPath("../cmd/cfg")
	viper.SetConfigType("yml")
	viper.SetConfigName(".wrkspc")
	errnie.Handles(viper.ReadInConfig())

	// Overwrite the stage to `test` so we are in the correct context.
	viper.Set("wrkspc.stage", "test")

	return &ConfigContext{}
}

func TestNew(t *testing.T) {
	NewConfigContext()

	Convey("Given a new instance", t, func() {
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
