package datura

import (
	"errors"

	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

type Radix struct {
	tree *iradix.Tree
	pool *twoface.Pool
}

func NewRadix() *Radix {
	pool := twoface.NewPool(twoface.NewContext())
	pool.Run()

	return &Radix{
		tree: iradix.New(),
		pool: pool,
	}
}

func (store *Radix) PoolSize() int {
	return store.pool.Size()
}

type readJob struct {
	store *Radix
	p     []byte
}

func (job readJob) Do() {
	it := job.store.tree.Root().Iterator()
	it.SeekLowerBound(spd.Unmarshal(job.p).Payload())

	// I honestly don't fully get what is going on in this for loop...
	for key, blob, ok := it.Next(); ok; key, blob, ok = it.Next() {
		_ = key
		// Hmm, this would be even faster as a channel. Let's do that.
		job.p = make([]byte, len(blob.([]byte)))
		copy(job.p, blob.([]byte))
	}
}

func (store *Radix) Read(p []byte) (n int, err error) {
	store.pool.Do(readJob{
		store: store,
		p:     p,
	})

	return len(p), nil
}

type writeJob struct {
	store *Radix
	p     []byte
}

func (job writeJob) Do() {
	var ok bool

	if _, _, ok = job.store.tree.Insert(
		[]byte(spd.Unmarshal(job.p).Prefix()), job.p,
	); !ok {
		errnie.Handles(errors.New("no write")).With(errnie.NOOP)
	}
}

func (store *Radix) Write(p []byte) (n int, err error) {
	store.pool.Do(writeJob{
		store: store,
		p:     p,
	})

	return len(p), nil
}
