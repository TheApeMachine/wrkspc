package errnie

import "strings"

/*
ErrorType adds a strongly typed context to the error.
*/
type ErrorType uint

const (
	// NIL represents no error was generated.
	NIL ErrorType = iota
	// VALIDATION represents an error during validation of a value.
	VALIDATION
	// UNKNOWN represents an error with unknown cause.
	UNKNOWN
)

/*
Error is a custom wrapping around Go errors keeping errnie errors
fully compatible with conventions.
*/
type Error struct {
	error
	Type ErrorType
	Msg  string
	err  error
}

/*
NewError instantiates a new errnie Error and returns a pointer to it.
*/
func NewError(err error) *Error {
	return &Error{
		Type: getType(err),
		Msg:  err.Error(),
		err:  err,
	}
}

/*
Error implements the Go error interface by returning the error message.
*/
func (err *Error) Error() string {
	return err.Msg
}

/*
getType takes the first segment of the error message and uses it to
perform a lookup for the strong error type.
*/
func getType(err error) ErrorType {
	switch strings.Split(err.Error(), " ")[0] {
	case "":
		return NIL
	case "invalid":
		return VALIDATION
	default:
		return UNKNOWN
	}
}
