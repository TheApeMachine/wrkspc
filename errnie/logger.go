package errnie

type LogLevel uint

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
)

/*
Logger is an interface objects can implement if they want to act as
log output channels for errnie.
*/
type Logger interface {
	Error(...any)
	Warning(...any)
	Info(...any)
	Debug(...any)
}

/*
Warns is syntactic sugar to call the Warning method on
a Logger interface.
*/
func Warns(msgs ...any) {
	sendOut(WARNING, msgs...)
}

/*
Informs is syntactic sugar to call the Info method on
a Logger interface.
*/
func Informs(msgs ...any) {
	sendOut(INFO, msgs...)
}

/*
Debugs is syntactic sugar to call the Debug method on
a Logger interface.
*/
func Debugs(msgs ...any) {
	sendOut(DEBUG, msgs...)
}
