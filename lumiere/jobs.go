package lumiere

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
RenderJob turns styled structures into bytes that can be converted
to an output format. In most cases this would be SVG.
*/
type RenderJob struct {
	ctx     twoface.Context
	element Element
}

/*
Do the render job continuously to make the widget able to respond to
state changes, and handle cacncellation via the context.
*/
func (job *RenderJob) Do() errnie.Error {
	for {
		select {
		case <-job.ctx.TTL():
			return errnie.NewError(nil)
		default:
			job.ctx.Write(spd.NewCached(
				"lumiere",
				"render",
				"lumiere.wrkgrp.org",
				job.element.Render().String(),
			))
		}
	}
}
