package twoface

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Worker wraps a concurrent process that is able to process Job types
scheduled onto a Pool.
*/
type Worker interface {
	io.ReadWriteCloser
}

func NewWorker(workerType Worker) Worker {
	errnie.Trace()

	if workerType == nil {
		return NewProtoWorker(
			make(chan chan Job),
			make(chan Job),
			NewContext(nil),
		)
	}

	return workerType
}

type ProtoWorker struct {
	pool chan chan Job
	jobs chan Job
	ctx  Context
}

func NewProtoWorker(
	pool chan chan Job, jobs chan Job, ctx Context,
) *ProtoWorker {
	return &ProtoWorker{pool, jobs, ctx}
}

func (worker *ProtoWorker) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

	go func() {
		for {
			select {
			case <-worker.ctx.TTL():
				// Lifespan for this worker is over.
				worker.Close()
				return
			default:
				// Return the job channel to the worker pool.
				worker.pool <- worker.jobs

				// Pick up a new job if available.
				job := <-worker.jobs

				// Do the work by calling the interface method on the current
				// instance.
				job.Do()

			}
		}
	}()

	return
}

func (worker *ProtoWorker) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (worker *ProtoWorker) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}
