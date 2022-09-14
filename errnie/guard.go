package errnie

import (
	"log"
)

/*
Guard is a way to override any existing interrupt messages such as panics and
recover from them either silently or by invoking a custom recovery function.
This is very similar to Ruby's rescue method. Similar care should be taken as
well. This is made for situations where you have an actual solid recovery plan
and it could get stuck in an infinite crash loop if you are not careful. To make
this work you have to put in additional systems on top.
*/
type Guard struct {
	Err     error
	handler func()
}

/*
NewGuard constructs an instance of a guard that can live basically anywhere and sit
ready to either Check or Rescue. Rescue does what it says on the tin and recovers from
a panic situation, optionally with a bootstrap function to start the full recovery process.
Check I don't remember why it is there, but I will look it up in the original project this
all comes from.
*/
func NewGuard(handler func()) *Guard {
	return &Guard{
		handler: handler,
	}
}

/*
Check... I have no fucking clue what this does.
*/
func (guard *Guard) Check() {
	if guard.Err != nil {
		panic(guard.Err)
	}
}

/*
Rescue a method from errors and panics. The way to set this up is to make a
deferred call to a previously instantiated Guard. By deferring the Rescue call
no matter what happens the Rescue method will be called and it uses Go's built in
recover statement to circumvent the panic.
*/
func (guard *Guard) Rescue() func() {
	return func() {
		// Perform the recovery or check the cached error value in the guard instance.
		// I think it is becoming more clear, seemingly I allow you to convert an error
		// into a panic to recover from. No idea why.
		if r := recover(); r != nil || guard.Err != nil {
			if guard.handler == nil {
				log.Printf("%v:%v\n", r, guard.Err)
				return
			}

			guard.handler()
		}
	}
}
