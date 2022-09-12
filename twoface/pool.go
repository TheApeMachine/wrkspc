package twoface

import (
	"time"
)

/*
Pool is a set of Worker types, each running their own (pre-warmed) goroutine.
Any object that implements the Job interface is able to schedule work on the
worker pool, which keeps the amount of goroutines in check, while still being
able to benefit from high concurrency in all kinds of scenarios.
*/
type Pool struct {
	maxWorkers    int
	disposer      *Context
	workers       chan chan Job
	jobs          chan Job
	handles       []*Worker
	scaleInterval time.Duration
	scaleRate     int
	stats         int64
	loadCount     int
	overload      bool
}

/*
NewPool instantiates a worker pool with bound size of maxWorkers, taking in a
Context type to be able to cleanly cancel all of the sub processes it starts.
*/
func NewPool(disposer *Context) *Pool {
	return &Pool{
		maxWorkers:    0,
		disposer:      disposer,
		workers:       make(chan chan Job),
		jobs:          make(chan Job),
		handles:       make([]*Worker, 0),
		scaleInterval: 10,
		scaleRate:     10,
		loadCount:     0,
		overload:      false,
	}
}

/*
Do is the entry point for new jobs that want to be scheduled onto the worker pool.
*/
func (pool *Pool) Do(jobType Job) {
	// The jobs channel is buffered to prevent the program from blocking if all
	// workers are currently busy.
	pool.jobs <- NewJob(jobType)
}

/*
Run the workers, after creating and assigning them to the pool.
*/
func (pool *Pool) Run() {
	// Periodically  we will evaluate the performance of the worker
	// pool and potentially grow or shrink it.
	ticker := time.NewTicker(pool.scaleInterval * time.Millisecond)

	go func() {
		for {
			select {
			case <-pool.disposer.Done():
				return
			case <-ticker.C:
				// Check if we are overloaded or not.
				pool.checkLoad()

				// Grow the pool based on the current load.
				if !pool.grow() {
					// We only have to potentially shrink the pool if
					// we didn't just grow it.
					pool.shrink()
				}
			}
		}
	}()

	// Start the job scheduling process.
	go pool.dispatch()
}

func (pool *Pool) dispatch() {
	// Make sure that we cleanly close the channels if our dispatcher
	// returns for whatever reason.
	defer close(pool.jobs)
	defer close(pool.workers)

	for {
		select {
		case job := <-pool.jobs:
			// A new job was received from the jobs queue, get the first available
			// worker from the pool once ready.
			jobChannel := <-pool.workers
			// Then send the job to the worker for processing.
			jobChannel <- job
		case <-pool.disposer.Done():
			// The disposer was triggered, clean up, and bail out.
			for _, worker := range pool.handles {
				worker.Stop()
			}

			return
		}
	}
}
