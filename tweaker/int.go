package tweaker

/*
GetInt returns a value by the key that is passed in and
converts it to an int type.

This is a public method that proxies the call to the internal
method bound to the ambient context.
*/
func GetInt(key string) int { return cfg.getInt(key) }

/*
getInt returns a value by the key that is passed in and
converts it to an int type.

This is an internal method bound to the ambient context.
*/
func (cfg *Config) getInt(key string) int {
	return cfg.v.GetInt(cfg.withStage(key))
}
