package tweaker

import (
	"github.com/wrk-grp/errnie"
)

/*
GetIdentity returns a value that provides an identity for the
service that is currently running.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetIdentity() []byte {
	errnie.Trace()
	return cfg.getIdentity()
}

/*
getIdentity  returns a value that provides an identity for the
service that is currently running.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getIdentity() []byte {
	errnie.Trace()
	identity := cfg.stage() + "." + cfg.program()
	return []byte(identity)
}
