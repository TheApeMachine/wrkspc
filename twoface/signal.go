package twoface

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Signal is a wrapper around terminal interrupt signals that are fired
when typing ctrl+c for example.
*/
type Signal struct {
	stop       chan struct{}
	interrupts []os.Signal
	ctxs       []Context
}

/*
NewSignal sets up an interrupt signal handler.
Pass in a splat on contexts for graceful shutdown.
*/
func NewSignal(ctxs ...Context) Signal {
	return Signal{
		stop:       make(chan struct{}),
		interrupts: []os.Signal{os.Interrupt, syscall.SIGTERM},
		ctxs:       ctxs,
	}
}

func (sig Signal) Run() chan struct{} {
	c := make(chan os.Signal, 2)
	signal.Notify(c, sig.interrupts...)

	go func() {
		<-c
		close(sig.stop)
		<-c

		for _, ctx := range sig.ctxs {
			ctx.cancel()
			errnie.Logs("context cancelled").With(errnie.INFO)
		}

		os.Exit(1)
	}()

	return sig.stop
}
