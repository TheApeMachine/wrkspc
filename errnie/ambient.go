package errnie

import (
	slacker "github.com/slack-go/slack"
	"github.com/spf13/viper"
)

var ambctx *AmbientContext

func init() {
	ambctx = New()
}

/*
AmbientContext is globally accessible throughout the entire program to reduce friction when it
comes to an object used in many places.
*/
type AmbientContext struct {
	tracer  *Tracer
	loggers []LogChannel
	msgs    []interface{}
	OK      bool
	ERR     *Error
}

/*
New constructs the AmbientContext such that is becomes accessible.
*/
func New() *AmbientContext {
	program := viper.GetString("program")

	ambctx := new(AmbientContext)
	ambctx.tracer = NewTracer()
	ambctx.loggers = []LogChannel{
		NewLogger(&ConsoleLogger{}),
		NewLogger(&SlackLogger{
			client: slacker.New(
				viper.GetString(program+".slack.token"), slacker.OptionDebug(true),
			),
		}),
	}
	ambctx.OK = true
	ambctx.ERR = nil

	return ambctx
}

/*
Runtime can be used to debug values under the hood of the runtime.
*/
func Runtime(interval int) { ambctx.tracer.Runtime(interval) }

/*
Traces the stack and outputs debugging information.
*/
func Traces(flags ...bool) { ambctx.tracer.Inspect(flags...) }

/*
Logs proxies the call onto the AmbientContext.
*/
func Logs(msgs ...interface{}) *AmbientContext { return ambctx.Logs(msgs...) }

/*
Logs is an ambient method that proxies to the configured LogChannel.
*/
func (ambctx *AmbientContext) Logs(msgs ...interface{}) *AmbientContext {
	ambctx.msgs = msgs
	return ambctx
}

/*
With is a chainable method that turns our code structure into something semantically pleasing.
Example:
  errnie.Logs(err).With(errnie.ERROR)
*/
func (ambctx *AmbientContext) With(logLevel LogLevel) *AmbientContext {
	var err *Error
	ambctx.OK = true

	switch logLevel {
	case ERROR:
		err = ambctx.loggers[0].Error(ambctx.msgs...)
		_ = ambctx.loggers[1].Error(ambctx.msgs...)
	case WARNING:
		err = ambctx.loggers[0].Warning(ambctx.msgs...)
		_ = ambctx.loggers[1].Warning(ambctx.msgs...)
	case INFO:
		ambctx.loggers[0].Info(ambctx.msgs...)
		ambctx.loggers[1].Info(ambctx.msgs...)
	case DEBUG:
		ambctx.loggers[0].Debug(ambctx.msgs...)
		ambctx.loggers[1].Debug(ambctx.msgs...)
	case INSPECT:
		ambctx.loggers[0].Inspect(ambctx.msgs...)
		ambctx.loggers[1].Inspect(ambctx.msgs...)
	}

	// Set OK to false when an error was indeed found. This is just a helper for the caller
	// to easily use in conditional statements. Same with ERR, but you can get the value that way.
	if err != nil {
		ambctx.OK = false
		ambctx.ERR = err
	}

	return ambctx
}
