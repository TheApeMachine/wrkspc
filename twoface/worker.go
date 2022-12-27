package twoface

import (
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Worker wraps a concurrent process that is able to process Job types
scheduled onto a Pool.
*/
type Worker struct {
	pool chan chan Job
	jobs chan Job
}

func NewWorker(pool chan chan Job, jobs chan Job) *Worker {
	errnie.Trace()
	return &Worker{pool, jobs}
}

func (worker *Worker) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (worker *Worker) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (worker *Worker) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}
