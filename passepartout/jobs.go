package passepartout

import (
	"sync"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
switchJob cancels out any parallel lookups that may be ongoing as soon
as the data has been found.
*/
func (manager *Manager) switchJob(
	p []byte, wgs []*sync.WaitGroup,
) (n int, err error) {
	errnie.Traces()
	for _, wg := range wgs {
		wg.Wait()

		if p != nil {
			// We have found our data, so we return.
			return len(p), nil
		}
	}

	return len(p), err
}

/*
ReadJob allows a lookup to be scheduled onto the pool of
pre-warmed worker routines.
*/
type ReadJob struct {
	store Store
	p     []byte
	wg    *sync.WaitGroup
}

/*
Do implements the twoface.Job interface.
*/
func (job ReadJob) Do() errnie.Error {
	errnie.Traces()
	defer job.wg.Done()
	job.store.Read(job.p)
	return errnie.NewError(nil)
}

/*
ManagerWriteJob makes it possible to schedule data write operations
onto the pre-warmed pool of worker routines.
*/
type ManagerWriteJob struct {
	store Store
	p     []byte
}

/*
Do implements the twoface.Job interface.
*/
func (job ManagerWriteJob) Do() errnie.Error {
	errnie.Traces()
	job.store.Write(job.p)
	return errnie.NewError(nil)
}
