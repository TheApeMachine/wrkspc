package passepartout

import (
	"sync"

	"github.com/theapemachine/wrkspc/twoface"
)

/*
Router builds dynamic pathways to data.
*/
type Router struct {
	ctx    *twoface.Context
	pool   *twoface.Pool
	routes *sync.Map
}

func NewRouter() *Router {
	ctx := twoface.NewContext()

	return &Router{
		ctx:    ctx,
		pool:   twoface.NewPool(ctx).Run(),
		routes: &sync.Map{},
	}
}

/*
Read implements io.Reader and is used to parse incoming data from
requests and map it to a route. If a route does not exist yet, it
will be dynamically created.
*/
func (router *Router) Read(p []byte) (n int, err error) {
	return
}

/*
Write implements io.Writer and is used to parse incoming data from
requests and map it to a route. If a route does not exist yet, it
will be dynamically created.
*/
func (router *Router) Write(p []byte) (n int, err error) {
	return
}

/*
PoolSize returns the current size of the autoscaling worker pool.
*/
func (router *Router) PoolSize() int {
	return router.pool.Size()
}
