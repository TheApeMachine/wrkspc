package errnie

import (
	"os"
)

var ctx *Context

func init() {
	ctx = New()
}

type Context struct {
	tracing   bool
	debugging bool
	loggers   []Logger
}

func New() *Context {
	return &Context{
		loggers: []Logger{NewConsoleLogger()},
	}
}

func Tracing(value bool) {
	ctx.tracing = value
}

func Debugging(value bool) {
	ctx.debugging = value
}

func Informs(msgs ...any) {
	sendOut(msgs...)
}

func Debugs(msgs ...any) {
	sendOut(msgs...)
}

func Handles(err error) *Error {
	if err == nil {
		return nil
	}

	out := NewError((err))
	sendOut(out)

	return out
}

func Kills(err error) {
	if err == nil {
		return
	}

	sendOut(err)
	os.Exit(1)
}

func sendOut(msgs ...any) {
	for _, logger := range ctx.loggers {
		logger.Debug(msgs...)
	}
}
