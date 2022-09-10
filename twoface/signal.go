package twoface

import (
	"os"
	"os/signal"
	"syscall"
)

/*
Signal is a wrapper around terminal interrupt signals that are fired
when typing ctrl+c for example.
*/
type Signal struct {
	stop       chan struct{}
	interrupts []os.Signal
}

/*
NewSignal sets up an interrupt signal handler.
*/
func NewSignal() Signal {
	return Signal{
		stop:       make(chan struct{}),
		interrupts: []os.Signal{os.Interrupt, syscall.SIGTERM},
	}
}

func (sig Signal) Run() chan struct{} {
	c := make(chan os.Signal, 2)
	signal.Notify(c, sig.interrupts...)

	go func() {
		<-c
		close(sig.stop)
		<-c
		os.Exit(1)
	}()

	return sig.stop
}
