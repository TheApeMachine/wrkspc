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
	level    int
	samples  int
	overload bool
	lower    bool
	pool     *Pool
	maxIdle  time.Duration
}

/*
NewScaler constructs a scaler which controls the size of a
worker pool dynamically, based on the consumption of
machine resources.
*/
func NewScaler(pool *Pool) *Scaler {
	errnie.Traces()

	return &Scaler{
		interval: 100,
		rate:     10,
		stats:    0,
		period:   0,
		level:    1,
		samples:  3,
		overload: false,
		lower:    false,
		pool:     pool,
		maxIdle:  1 * time.Second,
	}
}

func (scaler *Scaler) Run() {
	// Periodically  we will evaluate the performance of the worker
	// pool and potentially grow or shrink it.
	ticker := time.NewTicker(scaler.interval * time.Millisecond)

	go func() {
		for {
			select {
			case <-scaler.pool.ctx.TTL():
				return
			case <-ticker.C:
				scaler.load()

				// Start growing the worker pool when it is not
				// overloaded and there are actually jobs in the queue.
				if !scaler.overload && len(scaler.pool.jobs) > 0 {
					scaler.Grow()
				}

				// Start shrinking the worker pool when overloaded.
				if scaler.overload {
					scaler.Shrink()
				}
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
	scaler.period++

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

	// We should only evaluate if we have previously
	// collected statistics.
	if prev == 0 || scaler.stats == 0 {
		return
	}

	// Take the average worker runtime duration and store it back
	// into stats, so it is ready for the next iteration.
	scaler.stats = scaler.stats / int64(count)

	lower := scaler.stats > prev

	if scaler.lower != lower {
		scaler.period = 0

		if scaler.level > 1 {
			scaler.level--
		}
	}

	scaler.lower = lower

	if scaler.period >= scaler.samples {
		errnie.Debugs("period", scaler.pool.Size())
		scaler.period = 0

		if scaler.level < 3 {
			scaler.level++
		}

		if scaler.stats > prev {
			scaler.overload = true
			return
		}

		scaler.overload = false
		return
	}

	errnie.Debugs("load", "prev", prev, "stats", scaler.stats)
}

/*
grow the size of the worker pool by the rate we have defined. It is
better not to try to scale by single steps, as it is too slow, even
when the interval is set really low.
TODO: This can be further optimized by making the rate more dynamic.
*/
func (scaler *Scaler) Grow() {
	errnie.Traces()

	if !scaler.overload {
		for i := 0; i < scaler.rate*scaler.level; i++ {
			// Create a new worker and start its inner process, give it its own
			// disposer so we have granular control over the workers and we could
			// potentially dynamically resize the scaler later.
			scaler.pool.handles = append(scaler.pool.handles, NewWorker(
				len(scaler.pool.handles), scaler.pool.workers, NewContext(nil),
			).Start())
		}
	}
}

func (scaler *Scaler) Shrink() {
	errnie.Traces()

	if len(scaler.pool.handles) == 0 {
		return
	}

	if scaler.overload {
		for i := 0; i < scaler.rate; i++ {
			// Stop the worker, once it finishes its current job.
			scaler.drain(scaler.pool.handles[i], i)
		}

		return
	}

	// Drain any workers that are just sitting around idling.
	for idx, handle := range scaler.pool.handles {
		if time.Since(handle.lastUse) > scaler.maxIdle {
			scaler.drain(handle, idx)
		}
	}
}

func (scaler *Scaler) drain(worker *Worker, i int) {
	worker.Drain()
	hCnt := 0

	copy(scaler.pool.handles[i:], scaler.pool.handles[i+1:])
	hCnt = len(scaler.pool.handles) - 1
	scaler.pool.handles[hCnt] = nil
	scaler.pool.handles = scaler.pool.handles[:hCnt]
}
