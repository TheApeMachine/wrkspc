package errnie

import "os"

func Handles(err error) *Error {
	if err == nil {
		return nil
	}

	out := NewError(err)
	sendOut(ERROR, out)

	return out
}

func Kills(err error) {
	if err == nil {
		return
	}

	sendOut(ERROR, err)
	os.Exit(1)
}
