package tweaker

import (
	"fmt"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
GetString returns a value by the key that is passed in and
converts it to a string type.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetString(key string) string {
	errnie.Trace()
	return cfg.getString(key)
}

/*
getString returns a value by the key that is passed in and
converts it to a string type.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getString(key string) string {
	errnie.Trace()
	k := cfg.withStage(key)
	v := cfg.v.GetString(k)
	errnie.Debugs(fmt.Sprintf("tweaker.getString(%s) ->", k), v)
	return v
}
