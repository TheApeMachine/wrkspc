package tweaker

import "github.com/theapemachine/wrkspc/errnie"

/*
GetBool returns a value by the key that is passed in and
converts it to a boolean type.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetBool(key string) bool {
	errnie.Trace()
	return cfg.getBool(key)
}

/*
getBool returns a value by the key that is passed in and
converts it to a boolean type.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getBool(key string) bool {
	errnie.Trace()
	return cfg.v.GetBool(cfg.withStage(key))
}
