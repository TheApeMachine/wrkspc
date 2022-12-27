package tweaker

import "github.com/theapemachine/wrkspc/errnie"

/*
GetInt returns a value by the key that is passed in and
converts it to an int type.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetInt(key string) int {
	errnie.Trace()
	return cfg.getInt(key)
}

/*
getInt returns a value by the key that is passed in and
converts it to an int type.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getInt(key string) int {
	errnie.Trace()
	return cfg.v.GetInt(cfg.withStage(key))
}
