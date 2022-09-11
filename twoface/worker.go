package twoface

import (
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

type Worker struct {
	ID           int
	WorkerPool   chan chan Job
	JobChannel   chan Job
	disposer     Context
	working      bool
	idleCount    int
	lastDuration int64
	drain        bool
}

func NewWorker(
	ID int,
	workerPool chan chan Job,
	disposer Context,
) *Worker {
	return &Worker{
		ID:         ID,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		disposer:   disposer,
		working:    false,
		drain:      false,
	}
}

func (worker *Worker) Start() *Worker {
	go func() {
		defer close(worker.JobChannel)

		for {
			worker.WorkerPool <- worker.JobChannel

			select {
			case job := <-worker.JobChannel:
				worker.working = true
				worker.idleCount = 0
				t := time.Now()

				job.Do()

				worker.lastDuration = time.Since(t).Nanoseconds()
				worker.working = false

				if worker.drain {
					return
				}
			}
		}
	}()

	return worker
}

func (worker *Worker) Stop() {
	errnie.Traces()
	worker.disposer.cancel()
}

/*
Drain the worker, which means it will finish its current job first
before it will stop.
*/
func (worker *Worker) Drain() {
	errnie.Traces()
	worker.drain = true
}
