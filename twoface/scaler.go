package twoface

import (
	"math/rand"

	"github.com/theapemachine/wrkspc/errnie"
)

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
		if pool.stats > prev {
			pool.loadCount++

			if pool.loadCount < 3 {
				// We only want to scale down if there is a continuing
				// trend downwards.
				return
			}

			pool.loadCount = 0

			errnie.Logs(
				"overload", pool.maxWorkers,
			).With(errnie.WARNING)

			pool.overload = true
		}
	}
}

func (pool *Pool) grow() bool {
	errnie.Traces()

	if !pool.overload {
		for i := 0; i < pool.scaleRate; i++ {
			// Create a new worker and start its inner process, give it its own
			// disposer so we have granular control over the workers and we could
			// potentially dynamically resize the pool later.
			pool.handles = append(pool.handles, NewWorker(
				len(pool.handles), pool.workers, *NewContext(),
			).Start())

			pool.maxWorkers++
		}

		return true
	}

	return false
}

func (pool *Pool) shrink() {
	errnie.Traces()

	if pool.maxWorkers == 0 {
		return
	}

	if pool.overload && pool.maxWorkers > 1 {
		for i := 0; i < pool.scaleRate/2; i++ {
			// Pool is currently overloaded, start taking
			// out random workers.
			x := rand.Intn(pool.maxWorkers - 0)
			worker := pool.handles[x]

			// Stop the worker, once it finishes its current job.
			pool.drain(worker, x)
		}

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
			pool.drain(worker, i)
		}
	}
}

func (pool *Pool) drain(worker *Worker, i int) {
	worker.Drain()
	copy(pool.handles[i:], pool.handles[i+1:])
	pool.handles[len(pool.handles)-1] = nil
	pool.handles = pool.handles[:len(pool.handles)-1]
	pool.maxWorkers--
}
