package datura

import (
	"sync"

	"github.com/theapemachine/wrkspc/twoface"
)

/*
Manager act as a control structure to orchestrate a data lookup using
various data stores which act as canonical or caching resources.
*/
type Manager struct {
	stores []Store
	pool   *twoface.Pool
	ctx    *twoface.Context
}

/*
NewManager constructs a store manager capable of finding the shortest
distance and time to data.
*/
func NewManager() *Manager {
	ctx := twoface.NewContext()

	// Unfortunately for now, the order of stores matters, because of
	// the way we do a parallel lookup, and cancelling the S3 job if
	// the data was found in the Radix Tree.
	return &Manager{
		stores: []Store{NewRadix(), NewS3()},
		pool:   twoface.NewPool(ctx).Run(),
		ctx:    ctx,
	}
}

/*
ManagerReadJob allows a lookup to be scheduled onto the pool of
pre-warmed worker routines.
*/
type ManagerReadJob struct {
	store Store
	p     []byte
	wg    *sync.WaitGroup
}

/*
Do implements the twoface.Job interface.
*/
func (job ManagerReadJob) Do() {
	defer job.wg.Done()
	job.store.Read(job.p)
}

/*
switchJob cancels out any parallel lookups that may be ongoing as soon
as the data has been found.
*/
func (manager *Manager) switchJob(p []byte, wg []*sync.WaitGroup) (n int, err error) {
	// Both jobs (radix + S3) are running in goroutines, we wait for the fastest data
	// store to finish its lookup.
	wg[0].Wait()

	if p != nil {
		// We have found our data, so we return.
		return len(p), nil
	}

	// We have not found our data yet, wait for the S3 lookup to complete.
	wg[1].Wait()
	return len(p), err
}

/*
Read implements the io.Reader interface, and is used to perform a lookup using
one or more data stores.
*/
func (manager *Manager) Read(p []byte) (n int, err error) {
	// We need two separate wait groups to make a fall-through switch.
	wg := []*sync.WaitGroup{{}, {}}
	wg[0].Add(1) // Add one to the first wait group.
	wg[1].Add(1) // And one to the second wait group.

	// Loop over the attached data stores and start a new lookup
	// job for each one, passing in the wait groups with the
	// according index in the wait group array.
	for idx, store := range manager.stores {
		manager.pool.Do(ManagerReadJob{
			store: store,
			p:     p,
			wg:    wg[idx],
		})
	}

	// Reset p to be ready to receive result data from the lookup.
	p = nil

	// Call the logic switch to cancel parallel jobs when the results
	// are already in from another data store.
	return manager.switchJob(p, wg)
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
func (job ManagerWriteJob) Do() {
	job.store.Write(job.p)
}

/*
Write implements the io.Writer interface, and allows us to write data
into our data stores.
*/
func (manager *Manager) Write(p []byte) (n int, err error) {
	for _, store := range manager.stores {
		manager.pool.Do(ManagerWriteJob{
			store: store,
			p:     p,
		})
	}

	return len(p), err
}

/*
PoolSize completes the Employer interface implementation and returns the
current worker routine count, useful for debugging the auto-scaling pool.
*/
func (manager *Manager) PoolSize() int {
	return manager.pool.Size()
}
