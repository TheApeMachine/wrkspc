package errnie

/*
ctx holds our ambient context.
*/
var ctx *Context

/*
init loads up the ambient context so errnie can be called
from anywhere without instantiation.
*/
func init() {
	ctx = New()
}

/*
Context wraps all data and behavior errnie needs to act as
error handler, logger, and tracer.
*/
type Context struct {
	tracing     bool
	debugging   bool
	breakpoints bool
	loggers     []Logger
}

/*
New is a constructor to load up the ambient context with the
default values.
*/
func New() *Context {
	return &Context{
		loggers: []Logger{NewConsoleLogger()},
	}
}

/*
Tracing modifies the ambient context to turn tracing off or on.
*/
func Tracing(value bool) {
	ctx.tracing = value
}

/*
Debugging modifies the ambient context to turn debugging off or on.
*/
func Debugging(value bool) {
	ctx.debugging = value
}

/*
Breakpoints modifies the ambient context to turn breakpoints off or on.
*/
func Breakpoints(value bool) {
	ctx.breakpoints = value
}

/*
sendOut loops over all the log output channels that are active and
calls the currently activated logging method on them.
*/
func sendOut(level LogLevel, msgs ...any) {
	for _, logger := range ctx.loggers {
		switch level {
		case ERROR:
			logger.Error(msgs...)
		case WARNING:
			logger.Warning(msgs...)
		case INFO:
			logger.Info(msgs...)
		case DEBUG:
			logger.Debug(msgs...)
		}
	}
}
