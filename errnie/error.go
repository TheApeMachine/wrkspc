package errnie

/*
ErrorType adds a stronger context to Go error types.
*/
type ErrorType uint

const (
	// NIL represents the empty value for an error.
	NIL ErrorType = iota
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
	out := Error{Msg: err.Error()}

	switch out.Msg {
	case "":
		out.Type = NIL
	case "test error":
		out.Type = TEST
	}

	return out
}

func (err Error) Error() string {
	return err.Msg
}
