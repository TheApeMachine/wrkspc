package errnie

import (
	"strings"

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
	ec  ErrorContext
	et  LogLevel
	msg string
}

type Presenter struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

/*
NewError instantiates a new errnie Error and returns a pointer to it.
*/
func NewError(err error) *Error {
	if err == nil {
		return nil
	}

	split := strings.Split(err.Error(), "|")
	ec, et := getType(split[0])

	return &Error{
		ec:  ec,
		et:  et,
		msg: split[1],
	}
}

func (ee *Error) Read(p []byte) (n int, err error) {
	buf, err := json.Marshal(err)
	Handles(err)

	copy(p, buf)
	return len(p), err
}

func (ee *Error) Write(p []byte) (n int, err error) {
	Handles(json.Unmarshal(p, ee))
	return len(p), err
}

/*
Error implements the Go error interface by returning the error message.
*/
func (ee *Error) Error() string {
	return ee.msg
}

/*
Dump the error to the log output channels.
*/
func (ee *Error) Dump() {
	sendOut(ee.et, ee.msg)
}

/*
getType takes the first segment of the error message and uses it to
perform a lookup for the strong error type.
*/
func getType(err string) (ErrorContext, LogLevel) {
	split := strings.Split(err, " ")
	var ec ErrorContext
	var et LogLevel

	switch split[0] {
	case "":
		ec = NIL
	case "VALIDATION":
		ec = VALIDATION
	default:
		ec = UNKNOWN
	}

	switch split[1] {
	case "ERROR":
		et = ERROR
	case "WARNING":
		et = WARNING
	}

	return ec, et
}
