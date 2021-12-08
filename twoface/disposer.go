package twoface

import (
	"context"
)

/*
Disposer is a wrapper around cancellation contexts and part of
the concurrency `primitives` that are designed to be easy to use.
*/
type Disposer struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

/*
NewDisposer constructs a Disposer and returns a pointer reference
to it.
*/
func NewDisposer() *Disposer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Disposer{Ctx: ctx, Cancel: cancel}
}

/*
Cleanup triggers the disposer to send an empty struct to the Done
channel which will signal all the listeners to clean up after themselves.
*/
func (disposer *Disposer) Cleanup() {
	disposer.Cancel()
}

/*
Done returns the inner channel that signals an unreachable process
to start cleaning up and terminating.
*/
func (disposer *Disposer) Done() <-chan struct{} {
	return disposer.Ctx.Done()
}
