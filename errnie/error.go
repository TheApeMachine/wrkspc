package errnie

import (
	jsoniter "github.com/json-iterator/go"
)

/*
json provides us with a faster way to marshal and unmarshal
json data. This will replace the json package from the standard
library globally in any project that includes this package.
*/
var json = jsoniter.ConfigCompatibleWithStandardLibrary

/*
ErrorContext adds a strongly typed context to the error.
*/
type ErrorContext uint

const (
	// NIL represents no error was generated.
	NIL ErrorContext = iota
	// INTEGRITY represents a failure to verify the integrity of data.
	INTEGRITY
	// VALIDATION represents an error during validation of a value.
	VALIDATION
	// CONVERSION represents an error while converting a value type.
	CONVERSION
	// IOERROR represents a generic IO operation failure.
	IOERROR
	// UNKNOWN represents an error with unknown cause.
	UNKNOWN
)

/*
Error is a custom wrapping around Go errors keeping errnie errors
fully compatible with conventions.
*/
type Error struct {
	presenter Presenter
}

type Presenter struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

/*
NewError instantiates a new errnie Error and returns a pointer to it.
*/
func NewError(err error) error {
	return err
}

/*
Error implements the Go error interface by returning the error message.
*/
func (ee *Error) Error() string {
	return ee.presenter.Message
}
