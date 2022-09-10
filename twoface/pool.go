package twoface

import (
	"container/ring"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Pool is a set of Worker types, each running their own (pre-warmed) goroutine.
Any object that implements the Job interface is able to schedule work on the
worker pool, which keeps the amount of goroutines in check, while still being
able to benefit from high concurrency in all kinds of scenarios.
*/
type Pool struct {
	maxWorkers int
	disposer   *Context
	workers    chan chan Job
	jobs       chan Job
	handles    []Worker
	stats      *ring.Ring
}

/*
NewPool instantiates a worker pool with bound size of maxWorkers, taking in a
Context type to be able to cleanly cancel all of the sub processes it starts.
*/
func NewPool(disposer *Context) *Pool {
	return &Pool{
		maxWorkers: 1,
		disposer:   disposer,
		workers:    make(chan chan Job, 1),
		jobs:       make(chan Job, 1),
		handles:    make([]Worker, 1),
		stats:      ring.New(1),
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

func (pool *Pool) Grow(n int) {
	errnie.Traces()

	for i := pool.maxWorkers; i < pool.maxWorkers+n; i++ {
		pool.handles[i] = NewWorker(
			pool.workers, *NewContext(),
		).Start()
	}
}

func (pool *Pool) Shrink(n int) {
	errnie.Traces()

	for i := 0; i < n; i++ {
		l := len(pool.handles)
		pool.handles[i] = pool.handles[l-1]
		pool.handles = pool.handles[:l-1]
	}
}

/*
Run the workers, after creating and assigning them to the pool.
*/
func (pool *Pool) Run() {
	for i := 0; i < pool.maxWorkers; i++ {
		// Create a new worker and start its inner process, give it its own
		// disposer so we have granular control over the workers and we could
		// potentially dynamically resize the pool later.
		pool.handles[i] = NewWorker(
			pool.workers, *NewContext(),
		).Start()
	}

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
