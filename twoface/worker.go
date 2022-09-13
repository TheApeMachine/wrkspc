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
			// Return the job channel to the worker pool.
			worker.WorkerPool <- worker.JobChannel

			// Pick up a new job if available.
			job := <-worker.JobChannel

			// Reset the idle count to 0 always, because if this
			// reaches anything then 1, which is controller by the
			// pool scaler, the worker will be retired.
			worker.idleCount = 0
			worker.working = true
			t := time.Now()

			job.Do()

			worker.lastDuration = time.Since(t).Nanoseconds()
			worker.working = false

			if worker.drain {
				return
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
