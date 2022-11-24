package tweaker

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
cfg enables tweaker to be an ambient context.
*/
var cfg *Config

/*
init instantiates the ambient context before any other code
is executed.
*/
func init() {
	errnie.Trace()
	cfg = New()
}

/*
Config is the type that holds all data and bahavior for the
ambient context.

It wraps the viper package to provide a more convenient
interface when dealing with multiple stages/environments.
*/
type Config struct {
	v *viper.Viper
}

/*
New is a constructor method that helps initialize the
ambient context with its default values.
*/
func New() *Config {
	errnie.Trace()

	return &Config{
		v: viper.GetViper(),
	}
}

/*
program returns the value under the key with the same name
from the configuration file.
*/
func (cfg *Config) program() string {
	errnie.Trace()
	return cfg.v.GetString("program")
}

/*
stage returns the value under the key with the same name
from the configuration file.
*/
func (cfg *Config) stage() string {
	errnie.Trace()
	return cfg.v.GetString(cfg.program() + ".stage")
}

/*
withStage returns a key that drills down into the configuration
file up until the actual stages blocks.
*/
func (cfg *Config) withStage(key string) string {
	errnie.Trace()
	return strings.Join(
		[]string{cfg.program(), "stages", cfg.stage(), key}, ".",
	)
}
