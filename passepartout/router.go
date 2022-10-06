package passepartout

import (
	"io"
	"sync"

	"github.com/theapemachine/wrkspc/datura"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/sockpuppet"
	"github.com/theapemachine/wrkspc/spd"
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
	errnie.Traces()
	ctx := twoface.NewContext()

	return &Router{
		ctx:    ctx,
		pool:   twoface.NewPool(ctx).Run(),
		routes: &sync.Map{},
	}
}

/*
load a manager from the routes cache or make a new instance if the
requested object does not already exist.
*/
func (router *Router) load(
	role string, cached io.ReadWriter,
) io.ReadWriter {
	errnie.Traces()
	r, _ := router.routes.LoadOrStore(
		role, cached,
	)

	return r.(io.ReadWriter)
}

/*
readCached loads a manager from the routes cache, or instantiates it
first if it does not already exist. It then fills p with the return
value from the manager's distinct operation.
*/
func (router *Router) readCached(
	role string, cached io.ReadWriter, p []byte,
) {
	errnie.Traces()
	router.load(role, cached).Read(p)
}

/*
readCached loads a manager from the routes cache, or instantiates it
first if it does not already exist. It then writes p to the
manager's distinct operation.
*/
func (router *Router) writeCached(
	role string, cached io.ReadWriter, p []byte,
) {
	errnie.Traces()
	router.load(role, cached).Write(p)
}

/*
Read implements io.Reader and is used to parse incoming data from
requests and map it to a route. If a route does not exist yet, it
will be dynamically created.
*/
func (router *Router) Read(p []byte) (n int, err error) {
	errnie.Traces()
	dg := spd.Unmarshal(p)
	role, err := dg.Role()
	errnie.Handles(err)

	switch role {
	case "question":
		router.readCached(role, datura.NewS3(), p)
	case "service":
		router.readCached(role, sockpuppet.NewFastHTTPClient(), p)
	default:
		// Do whatever the prefix defines. This implements the read
		// part of endpoints as data.
		router.readCached(role, datura.NewS3(), p)
	}

	return
}

/*
Write implements io.Writer and is used to parse incoming data from
requests and map it to a route. If a route does not exist yet, it
will be dynamically created.
*/
func (router *Router) Write(p []byte) (n int, err error) {
	errnie.Traces()
	dg := spd.Unmarshal(p)
	role, err := dg.Role()
	errnie.Handles(err)

	switch role {
	case "datapoint":
		router.writeCached(role, datura.NewS3(), p)
	default:
		// Do whatever the prefix defines. This implements the write
		// part of endpoints as data.
		router.writeCached(role, datura.NewS3(), p)
	}

	return
}

/*
PoolSize returns the current size of the autoscaling worker pool.
*/
func (router *Router) PoolSize() int {
	errnie.Traces()
	return router.pool.Size()
}
