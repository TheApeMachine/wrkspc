package errnie

import (
	"fmt"
)

var ambctx AmbientContext

func init() {
	fmt.Println("errnie.init")
	ambctx = New()
}

/*
AmbientContext holds the handles to the objects exposed by errnie
as well as the data representing the current state.
*/
type AmbientContext struct {
	Status
	errors  []Error
	loggers []Logger
	logs    []interface{}
	tracer  Tracer
	Tracing bool
}

/*
New instantiates the AmbientContext so we can globally use errnie.
*/
func New() AmbientContext {
	fmt.Println("errnie.New")
	return AmbientContext{
		errors:  make([]Error, 0),
		loggers: []Logger{NewConsole()},
		tracer:  NewTracer(),
	}
}

func GetErrnie() Logger {
	return ambctx.loggers[0]
}

func Tracing(set bool) {
	ambctx.Tracing = set
}
