package errnie

import "strings"

type ErrorType uint

const (
	NIL ErrorType = iota
	UNK
)

type Error struct {
	error
	Type ErrorType
	Msg  string
	err  error
}

func NewError(err error) *Error {
	return &Error{
		Type: getType(err),
		Msg:  err.Error(),
		err:  err,
	}
}

func (err *Error) Error() string {
	return err.Msg
}

func getType(err error) ErrorType {
	switch strings.Split(err.Error(), " ")[0] {
	case "":
		return NIL
	default:
		return UNK
	}
}
