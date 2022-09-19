package datura

import (
	"bytes"
	"sync"

	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

var treeCache *iradix.Tree
var wgPool = &sync.Pool{
	New: func() interface{} {
		return &sync.WaitGroup{}
	},
}

type Radix struct {
	tree *iradix.Tree
	pool *twoface.Pool
}

func NewRadix() *Radix {
	pool := twoface.NewPool(twoface.NewContext())
	pool.Run()

	if treeCache == nil {
		treeCache = iradix.New()
	}

	return &Radix{
		tree: treeCache,
		pool: pool,
	}
}

func (store *Radix) PoolSize() int {
	return store.pool.Size()
}

type readJob struct {
	p  []byte
	wg *sync.WaitGroup
}

func (job readJob) Do() {
	defer job.wg.Done()

	it := treeCache.Root().Iterator()
	it.SeekPrefix(spd.Unmarshal(job.p).Payload())

	// I honestly don't fully get what is going on in this for loop...
	for key, blob, ok := it.Next(); ok; key, blob, ok = it.Next() {
		_ = key
		// Hmm, this would be even faster as a channel. Let's do that.
		buf := bytes.NewBuffer(job.p)
		buf.Truncate(0)
		bytes.NewBuffer(blob.([]byte)).WriteTo(buf)
	}
}

func (store *Radix) Read(p []byte) (n int, err error) {
	wg := wgPool.Get().(*sync.WaitGroup)
	wg.Add(1)

	store.pool.Do(readJob{
		p:  p,
		wg: wg,
	})

	wg.Wait()
	wgPool.Put(wg)
	return len(p), nil
}

type writeJob struct {
	p  []byte
	wg *sync.WaitGroup
}

func (job writeJob) Do() {
	defer job.wg.Done()
	treeCache, _, _ = treeCache.Insert(
		[]byte(spd.Unmarshal(job.p).Prefix()), job.p,
	)
}

func (store *Radix) Write(p []byte) (n int, err error) {
	wg := wgPool.Get().(*sync.WaitGroup)
	wg.Add(1)

	store.pool.Do(writeJob{
		p:  p,
		wg: wg,
	})

	wg.Wait()
	wgPool.Put(wg)
	return len(p), nil
}
