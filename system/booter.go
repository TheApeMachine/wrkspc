package system

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/gadget"
	"github.com/theapemachine/wrkspc/tweaker"
)

type Booter interface {
	Kick() chan error
}

func Boot(booterTypes ...Booter) chan error {
	errnie.Trace()
	out := make(chan error)

	for _, boot := range booterTypes {
		out <- <-boot.Kick()
	}

	return out
}

type SystemBooter struct{}

func (booter *SystemBooter) Kick() chan error {
	errnie.Tracing(tweaker.GetBool("errnie.trace"))
	errnie.Debugging(tweaker.GetBool("errnie.debug"))
	errnie.Breakpoints(tweaker.GetBool("errnie.break"))

	return gadget.NewPyroscope(
		tweaker.GetString("metrics.pyroscope.endpoint"),
	).Start()
}
