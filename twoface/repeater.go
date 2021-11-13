package twoface

/*
Repeater takes a simple higher order function and repeats it x amount of attempts to see
if it can receive a success.
*/
type Repeater struct {
	attempts      int
	retryStrategy RetryStrategy
}

/*
NewRepeater constructs a Repeater and returns a reference to it.
*/
func NewRepeater(attempts int, retryStrategy RetryStrategy) *Repeater {
	return &Repeater{
		attempts:      attempts,
		retryStrategy: retryStrategy,
	}
}

/*
Attempt proxies the work to the RetryStrategy.
*/
func (repeater *Repeater) Attempt(n int, try func() bool) int {
	return repeater.retryStrategy.Attempt(n, try)
}
