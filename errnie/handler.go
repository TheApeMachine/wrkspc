package errnie

import (
	"os"
	"runtime"
)

func Handles(err error) *Error {
	if out := NewError(err); out != nil {
		sendOut(ERROR, out.msg)
		Trap()
		return out
	}

	return nil
}

/*
Trap will execute a breakpoint trap if the relevant
configuration parameter is set to `true`.
*/
func Trap() {
	if ctx.breakpoints {
		runtime.Breakpoint()
	}
}

func Kills(err error) {
	if out := NewError(err); out != nil {
		sendOut(ERROR, out)
		os.Exit(1)
	}
}
