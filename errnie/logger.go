package errnie

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kind/pkg/log"
)

type Logger interface {
	Print(string, string, string, string)
	Error(string)
	Errorf(string, ...interface{})
	Warn(string)
	Warnf(string, ...interface{})
	V(log.Level) log.InfoLogger
}

type Log struct {
	Value string
}

func NewLog(val string) Log {
	return Log{
		Value: val,
	}
}

/*
Logs is a conveniance method to send values into the logging pipeline.
*/
func Logs(vals ...interface{}) Log {
	var builder strings.Builder

	for idx, val := range vals {
		switch v := val.(type) {
		case string:
			builder.WriteString(v)
		default:
			builder.WriteString(fmt.Sprintf("%v", v))
		}

		if idx < len(vals) {
			builder.WriteString(" ")
		}
	}

	return NewLog(builder.String())
}

/*
With is a chaining method that defines the follow on behavior to
apply to the error wrapped in the Handles method.
*/
func (log Log) With(op Op) Log {
	t, c, i := op()

	if ambctx.Debugging || t == " INFO  " {
		ambctx.loggers[0].Print(log.Value, t, c, i)
	}

	return log
}

/*
NewLogger converts a struct type into a Logger interface type.
*/
func NewLogger(loggerType Logger) Logger {
	return loggerType
}
