package twoface

import (
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Scaler controls the size of a worker pool and dynamically scales the
amount of worker routines, based on the current workload. This allows
us to redirect system resources to other pools when the load is low,
or allocate more resources when the load is high.
*/
type Scaler struct {
	interval time.Duration
	rate     int
	stats    int64
	period   int
	overload bool
	pool     *Pool
}

func NewScaler(pool *Pool) *Scaler {
	errnie.Traces()

	return &Scaler{
		interval: 1,
		rate:     20,
		stats:    0,
		period:   0,
		overload: false,
		pool:     pool,
	}
}

func (scaler *Scaler) Run() {
	// Periodically  we will evaluate the performance of the worker
	// pool and potentially grow or shrink it.
	ticker := time.NewTicker(scaler.interval * time.Millisecond)

	go func() {
		for {
			select {
			case <-scaler.pool.disposer.Done():
				return
			case <-ticker.C:
				scaler.load()
				scaler.Grow()
				scaler.Shrink()
			}
		}
	}()
}

/*
load determines if we see a degradation in the pool's performance
across a period (number of interations) and if so will set the overload
property of the scaler to true, indicating we should scale down.
*/
func (scaler *Scaler) load() {
	errnie.Traces()

	// Start with the default value, so we always have something
	// to return.
	scaler.overload = false

	// stats will contain the value from the previous iteration, which
	// we need to keep around to compare against.
	prev := scaler.stats
	scaler.stats = 0

	var count int
	var worker *Worker

	// Loop over all the current workers and sum their last runtime
	// durations. We can devide this by the size of the worker pool
	// to get a nice average.
	for _, worker = range scaler.pool.handles {
		if worker.lastDuration != 0 {
			scaler.stats += worker.lastDuration
			count++
		}
	}

	if count == 0 || scaler.stats == 0 {
		return
	}

	// Take the average worker runtime duration and store it back
	// into stats, so it is ready for the next iteration.
	scaler.stats = scaler.stats / int64(count)

	// We should only evaluate if we have previously
	// collected statistics.
	if prev == 0 {
		return
	}

	if scaler.stats > prev {
		// Our current average worker runtime duration is higher than
		// the iteration before. We want to see if this is a trend, so
		// for now just increase the period counter.
		scaler.period++

		if scaler.period < 3 {
			// We have not collected enough samples to determine
			// if there is a trend downwards.
			return
		}

		// A trend of performance degradation was detected. Reset the
		// period counter and indicate we are overloaded.
		scaler.period = 0
		scaler.overload = true

		return
	}
}

/*
grow the size of the worker pool by the rate we have defined. It is
better not to try to scale by single steps, as it is too slow, even
when the interval is set really low.
TODO: This can be further optimized by making the rate more dynamic.
*/
func (scaler *Scaler) Grow() bool {
	errnie.Traces()

	if !scaler.overload {
		for i := 0; i < scaler.rate; i++ {
			// Create a new worker and start its inner process, give it its own
			// disposer so we have granular control over the workers and we could
			// potentially dynamically resize the scaler later.
			scaler.pool.handles = append(scaler.pool.handles, NewWorker(
				len(scaler.pool.handles), scaler.pool.workers, *NewContext(),
			).Start())
		}

		return true
	}

	return false
}

func (scaler *Scaler) Shrink() {
	errnie.Traces()

	if len(scaler.pool.handles) == 0 {
		return
	}

	if scaler.overload {
		for i := 0; i < scaler.rate/2; i++ {
			// Stop the worker, once it finishes its current job.
			scaler.drain(scaler.pool.handles[i], i)
		}

		return
	}

	for i := 0; i < len(scaler.pool.workers); i++ {
		worker := scaler.pool.handles[i]

		if !worker.working {
			// The worker is currently not working, increase
			// the idleCount.
			worker.idleCount++
		}

		// We want to make sure the worker is actually idling. It will
		// reset idleCount every time it starts a job, so if we see
		// and idleCount of 2, we can be reasonably sure.
		if worker.idleCount >= 1 {
			// This worker ain't doing shit. Schedule for
			// death by shrinking.
			scaler.drain(worker, i)
		}
	}
}

func (scaler *Scaler) drain(worker *Worker, i int) {
	worker.Drain()
	copy(scaler.pool.handles[i:], scaler.pool.handles[i+1:])
	scaler.pool.handles[len(scaler.pool.handles)-1] = nil
	scaler.pool.handles = scaler.pool.handles[:len(scaler.pool.handles)-1]
}
