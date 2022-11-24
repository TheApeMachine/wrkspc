package twoface

import "github.com/theapemachine/wrkspc/errnie"

/*
Job is an interface any type can implement if they want to be able to use the
generics goroutine pool.
*/
type Job interface {
	Do() errnie.Error
}

/*
NewJob is a conveniance method to convert any incoming structured type to a
Job interface such that they can get onto the worker pools.
*/
func NewJob(jobType Job) Job {
	return jobType
}

/*
RetriableJob provides boilerplate for quickly building jobs that
retry based on a backoff delay strategy.
*/
type RetriableJob struct {
	ctx   Context
	fn    Job
	tries int
}

func NewRetriableJob(ctx Context, fn Job, tries int) Job {
	return NewJob(RetriableJob{
		ctx:   ctx,
		fn:    fn,
		tries: tries,
	})
}

/*
Do the job and retry x amount of times when needed.
*/
func (job RetriableJob) Do() errnie.Error {
	return NewRetrier(NewFibonacci(job.tries)).Do(job.fn)
}
