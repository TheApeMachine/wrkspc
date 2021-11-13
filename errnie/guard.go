package errnie

import (
	"fmt"
)

/*
Guard is a way to override any existing interrupt messages such as panics and
recover from them either silently or by invoking a custom recovery function.
*/
type Guard struct {
	Err     error
	logger  *Logger
	handler func()
}

/*
NewGuard constructs an instance of a guard that can live basically anywhere and sit
ready to either Check or Rescue.
*/
func NewGuard(handler func()) *Guard {
	return &Guard{
		logger:  NewLogger(ConsoleLogger{}),
		handler: handler,
	}
}

/*
Check is a method to force a guard to panic on an error and trigger it's own
recovery (Rescue) method. Be careful for infinite loops.
*/
func (guard *Guard) Check() {
	if guard.Err != nil {
		panic(guard.Err)
	}
}

/*
Rescue a method from errors and panics. If a handler was passed into the guard object
it will be the first thing that runs after recovery.
*/
func (guard *Guard) Rescue() func() {
	return func() {
		guard.recover()
	}
}

func (guard Guard) recover() {
	if r := recover(); r != nil || guard.Err != nil {
		guard.checkHandler(r)
	}
}

func (guard Guard) checkHandler(r interface{}) {
	if guard.handler == nil {
		guard.logger.Send(FATAL, (fmt.Sprintf("%v:%v", r, guard.Err)))
		return
	}

	guard.handler()
}
