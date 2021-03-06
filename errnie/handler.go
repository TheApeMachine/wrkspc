package errnie

import (
	"os"
)

/*
HandleFunc is a custom type we can use to define behavior we would like to
execute when the Handler has an error stored.
*/
type HandleFunc func(*Error)

var (
	// KILL exits the program.
	KILL HandleFunc = kill
	// NOOP does nothing, the error will still log.
	NOOP HandleFunc = noop
	// NOLO does nothing, the error will not log.
	NOLO HandleFunc = nolo
)

/*
Handler is the object used to contain the behavior we want to see based on the
type of error we have determined.
*/
type Handler struct {
	ERR *Error
	OK  bool
}

/*
Handles an error and passes back a pointer to the Handler object.
Really it is just a constructor, just one that also performs its
main action so what you get back is a `final` state, not an initialized one.
*/
func Handles(err error) *Handler {
	Traces()
	return &Handler{ERR: NewError(err)}
}

/*
With is a chainable method onto a Handler that allows behavior to be injected
that determines what to do when an error is detected.
*/
func (handler *Handler) With(fn HandleFunc) *Handler {
	Traces()
	handler.OK = true

	if len(handler.ERR.errs) == 0 || handler.ERR.errs[0] == nil {
		return handler
	}

	// Set OK to false when an error was indeed found. This is just a helper for the caller
	// to easily use in conditional statements.
	handler.OK = false

	fn(handler.ERR)

	return handler
}

/*
kill the program with exit.
*/
func kill(err *Error) {
	Traces()
	Logs(err.errs[0]).With(ERROR)
	os.Exit(1)
}

/*
noop can be used to log an error but ignore it otherwise.
*/
func noop(err *Error) {
	Traces()
	Logs(err.errs[0]).With(WARNING)
}

/*
nolo can be used to completely ignore an error and also not log it in any way.
*/
func nolo(err *Error) {
	Traces()
}
