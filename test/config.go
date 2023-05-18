package test

import (
	"github.com/spf13/viper"
	"github.com/wrk-grp/errnie"
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
