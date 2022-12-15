package errnie

import (
	"os"
	"runtime"
)

func Handles(err error) error {
	if out := NewError(err); out != nil {
		Errors(out)
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
		os.Exit(1)
	}
}
