package errnie

import (
	"errors"

	"github.com/davecgh/go-spew/spew"
)

/*
Op enumerates the operations errnie can perform on a error event.
*/
type Op func() (string, string, string)

var (
	// NOOP does not do anything.
	NOOP Op = writeLog(" NOOP  ", "HIGH", "ghost")
	// KILL exits the program with code 1.
	KILL Op = writeLog(" KILL  ", "HIGH", "skull")
	// SUCCESS ...
	SUCCESS Op = writeLog("SUCCESS", "HIGH", "thumu")
	// INFO ...
	INFO Op = writeLog(" INFO  ", "NORM", "badge")
	// DEBUG ...
	DEBUG Op = writeLog(" DEBUG ", "MUTE", "lbug")
	// WARNING ...
	WARNING Op = writeLog("WARNING", "NORM", "warn")
	// ERROR ...
	ERROR Op = writeLog(" ERROR ", "HIGH", "fire")
	// INSPECT ...
	INSPECT Op = writeInspect()
)

func writeLog(t, c, i string) func() (string, string, string) {
	return func() (string, string, string) {
		return t, c, i
	}
}

func writeInspect() func() (string, string, string) {
	return func() (string, string, string) {
		for _, log := range ambctx.logs {
			spew.Dump(log)
		}

		ambctx.logs = make([]interface{}, 0)
		return "", "", ""
	}
}

/*
Handles is a conveniance method that wraps a Go error value and
brings it into the errnie workflow.
*/
func Handles(err error) AmbientContext {
	if err == nil {
		return ambctx
	}

	ambctx.errors = append(ambctx.errors, NewError(err))
	// We always return the ambient context to keep methods chainable.
	return ambctx
}

func (ambctx AmbientContext) Return() func() func() {
	return NewGuard(nil).Rescue
}

/*
With is a chaining method that defines the follow on behavior to
apply to the error wrapped in the Handles method.
*/
func (ambctx AmbientContext) With(op Op) Error {
	defer op()

	if len(ambctx.errors) > 0 {
		return ambctx.errors[0]
	}

	return NewError(errors.New(""))
}
