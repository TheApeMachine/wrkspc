package errnie

import (
	"bytes"
	"encoding/gob"
)

/*
ErrorType defines canonical errors as a strong type so we can have a predicatable
value associated with them, which is a bnenefit in some cases over the string type
that sits underneath Go's builtin errors.
*/
type ErrorType uint

const (
	// ValidationError is returned when an inspected object does not contain all the
	// required values, or if any of those values are not set as expected.
	ValidationError ErrorType = iota
)

/*
Error is a thin wrapper around Go's builtin error type that adds strong typing and
some other neat functionality to it.
*/
type Error struct {
	errTypes []ErrorType
	errs     []error
}

/*
NewError wraps a builtin error into our custom type.
*/
func NewError(errs ...error) *Error {
	if len(errs) == 0 {
		return nil
	}

	return &Error{
		errTypes: make([]ErrorType, len(errs)),
		errs:     errs,
	}
}

/*
Encode the Error into a bytes.Buffer.
*/
func (err *Error) Encode() *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	encoder.Encode(*err)
	return buf
}

/*
First returns the first error of errs, since this type can hold multiple errors as one.
*/
func (err *Error) First() error {
	return err.errs[0]
}
