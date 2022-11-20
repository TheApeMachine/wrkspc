package tweaker

import (
	"strings"

	"github.com/spf13/viper"
)

var cfg *Config

func init() {
	cfg = New()
}

type Config struct {
	v *viper.Viper
}

func New() *Config {
	return &Config{
		v: viper.GetViper(),
	}
}

func (cfg *Config) program() string {
	return cfg.v.GetString("program")
}

func (cfg *Config) stage() string {
	return cfg.v.GetString(cfg.program() + ".stage")
}

func GetBool(key string) bool { return cfg.getBool(key) }

func (cfg *Config) getBool(key string) bool {
	return cfg.v.GetBool(strings.Join(
		[]string{cfg.program(), "stages", cfg.stage(), key}, ".",
	))
}
