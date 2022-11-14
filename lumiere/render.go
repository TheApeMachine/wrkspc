package lumiere

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Render is a presentable UI widget projected as an SVG.
*/
type Render struct {
	ctx  twoface.Context
	pool *twoface.Pool
}

/*
NewRender instantiates a projectable SVG UI widget.
*/
func NewRender(ctx twoface.Context) *Render {
	errnie.Traces()

	pool := twoface.NewPool(ctx)
	pool.Run()

	return &Render{
		ctx:  ctx,
		pool: pool,
	}
}

/*
Present the widget so it can be streamed over a web socket.
*/
func (render *Render) Present() []byte {
	render.pool.Do(
		&RenderJob{
			ctx:     render.ctx,
			element: NewButton(),
		},
	)
	return []byte{}
}
