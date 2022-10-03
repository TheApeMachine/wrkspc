package passepartout

import (
	"sync"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Manager is a generic means of passing on data and messages through
the internal service layers. It should be sufficient for most
use-cases.
*/
type Manager struct {
	ctx    *twoface.Context
	pool   *twoface.Pool
	stores []Store
}

/*
NewManager constructs a manager to move data around in a distributed
service architecture. Manager implements the Store interface itself,
meaning that managers are chainable.
*/
func NewManager(stores ...Store) *Manager {
	errnie.Traces()
	ctx := twoface.NewContext()

	return &Manager{
		ctx:    ctx,
		pool:   twoface.NewPool(ctx).Run(),
		stores: stores,
	}
}

/*
Read implements the io.Reader interface and is used to perform a
lookup operation on one or more stores.
Stores always compete to be the first to deliver the data and the
winner should cancel any other ongoing lookups.
*/
func (manager *Manager) Read(p []byte) (n int, err error) {
	errnie.Traces()
	var wgs []*sync.WaitGroup

	for idx, store := range manager.stores {
		var wg sync.WaitGroup
		wg.Add(1)
		wgs = append(wgs, &wg)

		manager.pool.Do(ReadJob{
			store: store,
			p:     p,
			wg:    wgs[idx],
		})
	}

	p = nil
	return manager.switchJob(p, wgs)
}

/*
Write implements the io.Writer interface, and allows us to write data
into our data stores.
*/
func (manager *Manager) Write(p []byte) (n int, err error) {
	errnie.Traces()
	for _, store := range manager.stores {
		manager.pool.Do(ManagerWriteJob{
			store: store,
			p:     p,
		})
	}

	return len(p), err
}

func (manager *Manager) PoolSize() int {
	errnie.Traces()
	return manager.pool.Size()
}
