package errnie

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
	ambctx := new(AmbientContext)
	ambctx.tracer = NewTracer()
	ambctx.loggers = []LogChannel{
		NewLogger(&ConsoleLogger{}),
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
func Traces(flags ...bool) string { return ambctx.tracer.Inspect(flags...) }

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
	ensureLogger()

	switch logLevel {
	case ERROR:
		for _, l := range ambctx.loggers {
			err = l.Error(ambctx.msgs...)
		}
	case WARNING:
		for _, l := range ambctx.loggers {
			err = l.Warning(ambctx.msgs...)
		}
	case INFO:
		for _, l := range ambctx.loggers {
			l.Info(ambctx.msgs...)
		}
	case DEBUG:
		for _, l := range ambctx.loggers {
			l.Debug(ambctx.msgs...)
		}
	case INSPECT:
		for _, l := range ambctx.loggers {
			l.Inspect(ambctx.msgs...)
		}
	}

	// Set OK to false when an error was indeed found. This is just a helper for the caller
	// to easily use in conditional statements. Same with ERR, but you can get the value that way.
	if err != nil {
		ambctx.OK = false
		ambctx.ERR = err
	}

	return ambctx
}
