package twoface

import (
	"math"
	"time"
)

/*
RetryStrategy is an interface you can implement if you don't like the ones that come out of the box.
*/
type RetryStrategy interface {
	Attempt(int, func() bool) int
}

/*
NewRetryStrategy constructs a RetryStrategy of the type that is passed in.
*/
func NewRetryStrategy(retryStrategyType RetryStrategy) RetryStrategy {
	return retryStrategyType
}

/*
Fibonacci uses the naturally increasing sequence as a backoff strategy.
*/
type Fibonacci struct {
	MaxTries int
}

/*
Attempt is the generic method on the interface that applies the retry strategy.
*/
func (strategy Fibonacci) Attempt(n int, try func() bool) int {
	// Enough is enough, bail.
	if n > strategy.MaxTries {
		return n
	}

	// Call the method passed into the retry strategy.
	if try() {
		// Success, bail.
		return 0
	}

	// Delay value is the positional value of fibonacci sequence.
	time.Sleep(time.Duration(n) * time.Second)

	// Recursion.
	return strategy.Attempt(
		int(
			math.Round((math.Pow(
				math.Phi, float64(n),
			)+math.Pow(
				math.Phi-1, float64(n),
			))/math.Sqrt(5)),
		), try, // Don't forget to pass in the operating function again.
	)
}
