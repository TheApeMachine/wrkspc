package twoface

import (
	"math"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Retrier is an interface that can be implemented by any object that wants to
schedule itself onto a worker pool and be retried under certain conditions.
*/
type Retrier interface {
	Do(Job) error
}

func NewRetrier(retrierType Retrier) Retrier {
	return retrierType
}

/*
Fibonacci is a RetryStategy that retries a function n times with a Fibonacci
interval in seconds between retries.
*/
type Fibonacci struct {
	max int
	n   int
}

func NewFibonacci(max int) Retrier {
	return NewRetrier(&Fibonacci{
		max: max,
		n:   0,
	})
}

func (strategy *Fibonacci) Do(fn Job) error {
	errnie.Trace()

	// We have reached the maximum number of retries.
	// Bail.
	if strategy.n > strategy.max {
		return errnie.NewError(nil)
	}

	// Error, retry.
	if err := fn.Do(); err != nil {
		// Backoff delay time by using Fibonacci sequence.
		strategy.n = int(
			math.Round((math.Pow(
				math.Phi, float64(strategy.n),
			) + math.Pow(
				math.Phi-1, float64(strategy.n),
			)) / math.Sqrt(5)),
		)

		// Wait for the next retry.
		time.Sleep(time.Duration(strategy.n) * time.Second)
		strategy.Do(fn)
	}

	return errnie.NewError(nil)
}
