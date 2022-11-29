package errnie

import "os"

func Handles(err error) *Error {
	if out := NewError(err); out != nil {
		sendOut(ERROR, out)
		return out
	}

	return nil
}

func Kills(err error) {
	if out := NewError(err); out != nil {
		sendOut(ERROR, out)
		os.Exit(1)
	}
}
