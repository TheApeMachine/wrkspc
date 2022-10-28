package passepartout

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

type Hub struct {
	ctx  *twoface.Context
	pool *twoface.Pool
}

func NewHub(ctx *twoface.Context) *Hub {
	errnie.Traces()

	return &Hub{
		ctx:  ctx,
		pool: twoface.NewPool(ctx).Run(),
	}
}

func (hub *Hub) Read(p []byte) (n int, err error) {
	return
}

func (hub *Hub) Write(p []byte) (n int, err error) {
	return
}

func (hub *Hub) PoolSize() int {
	errnie.Traces()
	return hub.pool.Size()
}
