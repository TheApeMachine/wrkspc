package tweaker

import (
	"fmt"

	"github.com/wrk-grp/errnie"
)

/*
GetSlice returns a value by the key that is passed in and
converts it to a string slice type.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetSlice(key string) []string {
	errnie.Trace()
	return cfg.getSlice(key)
}

/*
getString returns a value by the key that is passed in and
converts it to a string slice type.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getSlice(key string) []string {
	errnie.Trace()
	k := cfg.withStage(key)
	v := cfg.v.GetStringSlice(k)
	errnie.Debugs(fmt.Sprintf("tweaker.getSlice(%s) ->", k), v)
	return v
}
