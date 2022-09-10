package twoface

/*
Job is an interface any type can implement if they want to be able to use the
generics goroutine pool.
*/
type Job interface {
	Do()
}

/*
NewJob is a conveniance method to convert any incoming structured type to a
Job interface such that they can get onto the worker pools.
*/
func NewJob(jobType Job) Job {
	return jobType
}
