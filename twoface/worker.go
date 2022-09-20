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
	lastUse      time.Time
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
		lastUse:    time.Now(),
		drain:      false,
	}
}

/*
Start the worker to be ready to accept jobs from the job queue.
*/
func (worker *Worker) Start() *Worker {
	go func() {
		defer close(worker.JobChannel)

		for {
			// Return the job channel to the worker pool.
			worker.WorkerPool <- worker.JobChannel

			// Pick up a new job if available.
			job := <-worker.JobChannel

			// Keep track of the time before the work starts, with a
			// secondary benefit of helping to determine if the worker
			// is idle for a significant amount of time later on.
			worker.lastUse = time.Now()

			// Do the work by calling the interface method on the current
			// instance.
			job.Do()

			// Store the duration of the job load so it can later be used to
			// determine if the worker pool is overloaded.
			worker.lastDuration = time.Since(worker.lastUse).Nanoseconds()

			// This worker is about to get retired in a pool schrink.
			if worker.drain {
				return
			}
		}
	}()

	return worker
}

/*
Drain the worker, which means it will finish its current job first
before it will stop.
*/
func (worker *Worker) Drain() {
	errnie.Traces()
	worker.drain = true
}
