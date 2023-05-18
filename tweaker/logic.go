package tweaker

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/wrk-grp/errnie"
)

/*
GetLogic allows for more complex values in the configuration that require some
for of logic to be interpreted correctly.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetLogic(key string) int { return cfg.getLogic(key) }

/*
GetLogic allows for more complex values in the configuration that require some
for of logic to be interpreted correctly.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getLogic(key string) int {
	raw := strings.Split(cfg.v.GetString(cfg.withStage(key)), "*")

	switch raw[0] {
	case "cores":
		return runtime.NumCPU()
	case "threads":
		return runtime.NumCPU() * 2
	default:
		i, err := strconv.Atoi(raw[0])
		errnie.Handles(err)
		return i
	}
}
