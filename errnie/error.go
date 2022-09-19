package errnie

/*
ErrorType adds a stronger context to Go error types.
*/
type ErrorType uint

const (
	// NIL represents the empty value for an error.
	NIL ErrorType = iota
	// NOK represents a Not OK state.
	NOK
	// TEST represents a test context which should be ignored.
	TEST
)

/*
Error is a thin wrapper around Go errors.
*/
type Error struct {
	Msg  string
	Type ErrorType
}

/*
NewError constructs a new errnie Error type.
*/
func NewError(err error) Error {
	if err == nil {
		return Error{Type: NIL}
	}

	switch err.Error() {
	case "":
		return Error{Type: NIL}
	default:
		return Error{Type: NOK, Msg: err.Error()}
	}
}

func (err Error) Error() string {
	return err.Msg
}
