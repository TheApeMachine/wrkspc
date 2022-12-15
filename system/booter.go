package system

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/gadget"
	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/theapemachine/wrkspc/twoface"
)

type Booter interface {
	Kick() chan error
}

func Boot(booterTypes ...Booter) chan error {
	errnie.Trace()
	out := make(chan error)

	go func() {
		defer close(out)

		for _, boot := range booterTypes {
			out <- <-boot.Kick()
		}
	}()

	return out
}

type SystemBooter struct {
	Ctx *twoface.Context
	err error
}

func (booter *SystemBooter) Kick() chan error {
	errnie.Tracing(tweaker.GetBool("errnie.trace"))
	errnie.Debugging(tweaker.GetBool("errnie.debug"))
	errnie.Breakpoints(tweaker.GetBool("errnie.break"))

	errnie.Trace()
	out := make(chan error)

	go func() {
		defer close(out)
		errnie.Informs("system booting...")

		out <- gadget.NewPyroscope(
			tweaker.GetString("metrics.pyroscope.endpoint"),
		).Start()
	}()

	return out
}
