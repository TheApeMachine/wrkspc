package twoface

import (
	"fmt"
	"math/rand"
	"time"

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
	handles    []*Worker
	stats      int64
	overload   bool
}

/*
NewPool instantiates a worker pool with bound size of maxWorkers, taking in a
Context type to be able to cleanly cancel all of the sub processes it starts.
*/
func NewPool(disposer *Context) *Pool {
	return &Pool{
		maxWorkers: 0,
		disposer:   disposer,
		workers:    make(chan chan Job),
		jobs:       make(chan Job),
		handles:    make([]*Worker, 0),
		overload:   false,
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

func (pool *Pool) checkLoad() {
	pool.overload = false

	if pool.maxWorkers <= 1 {
		return
	}

	var count int64
	prev := pool.stats

	for _, worker := range pool.handles {
		if worker.lastDuration != 0 {
			pool.stats += worker.lastDuration
			count++
		}
	}

	if count == 0 {
		return
	}

	// Get the average duration of the jobs.
	pool.stats = pool.stats / count

	// We should only evaluate if we have previously
	// collected statistics.
	if prev > 0 {
		load := prev - pool.stats

		if load < 0 {
			errnie.Logs(fmt.Sprintf("pool overload: %d", load)).With(errnie.WARNING)
			pool.overload = true
		}
	}
}

func (pool *Pool) grow() bool {
	errnie.Traces()

	if !pool.overload {
		// Create a new worker and start its inner process, give it its own
		// disposer so we have granular control over the workers and we could
		// potentially dynamically resize the pool later.
		pool.handles = append(pool.handles, NewWorker(
			len(pool.handles), pool.workers, *NewContext(),
		).Start())

		pool.maxWorkers++

		errnie.Logs(
			fmt.Sprintf("new worker added (%d)", len(pool.handles)),
		).With(errnie.DEBUG)

		return true
	}

	return false
}

func (pool *Pool) shrink() {
	errnie.Traces()

	if pool.maxWorkers <= 1 {
		return
	}

	if pool.overload {
		// Pool is currently overloaded, start taking
		// out random workers.
		x := rand.Intn(pool.maxWorkers - 0)
		worker := pool.handles[x]

		// Stop the worker, once it finishes its current job.
		worker.Drain()
		copy(pool.handles[x:], pool.handles[x+1:])
		pool.handles[len(pool.handles)-1] = nil
		pool.handles = pool.handles[:len(pool.handles)-1]
		pool.maxWorkers--

		return
	}

	for i := 0; i < pool.maxWorkers; i++ {
		worker := pool.handles[i]

		if !worker.working {
			// The worker is currently not working, increase
			// the idleCount.
			worker.idleCount++
		}

		if worker.idleCount >= 3 {
			// This worker ain't doing shit. Schedule for
			// death by shrinking.
			worker.Drain()
			copy(pool.handles[i:], pool.handles[i+1:])
			pool.handles[len(pool.handles)-1] = nil
			pool.handles = pool.handles[:len(pool.handles)-1]
			pool.maxWorkers--
		}
	}
}

/*
Run the workers, after creating and assigning them to the pool.
*/
func (pool *Pool) Run() {
	// Periodically  we will evaluate the performance of the worker
	// pool and potentially grow or shrink it.
	ticker := time.NewTicker(300 * time.Millisecond)

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
