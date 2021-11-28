package errnie

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/berrt"
)

/*
LogLevel defines the criticality of the logged item.
*/
type LogLevel uint

const (
	// PANIC is a severe program crash.
	PANIC LogLevel = iota
	// FATAL is program crash.
	FATAL
	// CRITICAL is likely soon program crash.
	CRITICAL
	// ERROR is something went wrong.
	ERROR
	// WARNING is something that should be noticed.
	WARNING
	// INFO stands out as some relevant information.
	INFO
	// DEBUG provides context while investigating or evaluating the program.
	DEBUG
	// TRACE provides a look at values under the hood.
	TRACE
	// INSPECT dumps a deep analysis of the structure of objects.
	INSPECT
)

/*
LogChannel is an interface that allows logs to be sent to various solutions.
*/
type LogChannel interface {
	Error(...interface{}) *Error
	Warning(...interface{}) *Error
	Info(...interface{})
	Debug(...interface{})
	Inspect(...interface{})
}

/*
NewLogger is a constructor of LogChannel types.
*/
func NewLogger(logChannel LogChannel) LogChannel {
	return logChannel
}

/*
ConsoleLogger outputs to the terminal.
*/
type ConsoleLogger struct {
}

/*
Error logs the line with an error indicator and converts the slice of interfaces to Go's
strongly typed errors.
*/
func (logger ConsoleLogger) Error(events ...interface{}) *Error {
	if len(events) == 0 {
		return nil
	}

	var errs []error

	for _, err := range events {
		if err == nil {
			break
		}

		errs = append(errs, err.(error))
	}

	if len(errs) == 0 {
		return nil
	}

	fmt.Printf(
		"%s %s\n",
		berrt.NewLabel(" ERROR ").ToString(),
		berrt.NewText(fmt.Sprintf("%v", events...)).ToString(),
	)

	// Since internally we're a bit deeper in the stack, pass in an extra true to make
	// sure the tracing process jumps back far enough to reach the actual relevant code.
	Traces(true, true)

	return NewError(errs...)
}

/*
Warning logs the line with an error indicator and converts the slice of interfaces to Go's
strongly typed errors.
*/
func (logger ConsoleLogger) Warning(events ...interface{}) *Error {
	if len(events) == 0 {
		return nil
	}

	var errs []error

	for _, err := range events {
		if err == nil {
			break
		}

		errs = append(errs, err.(error))
	}

	if len(errs) == 0 {
		return nil
	}

	fmt.Printf(
		"%s %s\n",
		berrt.NewLabel("WARNING").ToString(),
		berrt.NewText(fmt.Sprintf("%v", events...)).ToString(),
	)

	// Since internally we're a bit deeper in the stack, pass in an extra true to make
	// sure the tracing process jumps back far enough to reach the actual relevant code.
	Traces(true, true)

	return NewError(errs...)
}

/*
Info is a helper output for informing the user.
*/
func (logger ConsoleLogger) Info(events ...interface{}) {
	if len(events) == 0 {
		return
	}

	var tmplt strings.Builder

	for range events {
		tmplt.WriteString("%v")
	}

	fmt.Printf(
		"%s %s\n",
		berrt.NewLabel(" INFO  ").ToString(),
		berrt.NewText(fmt.Sprintf(tmplt.String(), events...)).ToString(),
	)
}

/*
Debug is a helper output for development or troubleshooting.
*/
func (logger ConsoleLogger) Debug(events ...interface{}) {
	if len(events) == 0 || !viper.GetViper().GetBool("wrkspc.errnie.debug") {
		return
	}

	fmt.Printf(
		"%s %s\n",
		berrt.NewLabel(" DEBUG ").ToString(),
		berrt.NewText(fmt.Sprintf("%v", events...)).ToString(),
	)
}

/*
Inspect is a helper output for development or troubleshooting.
*/
func (logger ConsoleLogger) Inspect(events ...interface{}) {
	if len(events) == 0 {
		return
	}

	fmt.Printf(
		"%s %s\n",
		berrt.NewLabel("INSPECT").ToString(),
		berrt.NewText(spew.Sdump(events...)).ToString(),
	)
}
