package lumiere

import "github.com/theapemachine/wrkspc/twoface"

/*
RenderJob turns styled structures into bytes that can be converted
to an output format. In most cases this would be SVG.
*/
type RenderJob struct {
	ctx     *twoface.Context
	element Element
}

/*
Do the render job continuously to make the widget able to respond to
state changes, and handle cacncellation via the context.
*/
func (job *RenderJob) Do() {
	for {
		select {
		case <-job.ctx.Done():
			return
		default:
			// TODO: Probably better to put this onto some kind of
			// render queue that transports over Cap 'n Proto.
			job.ctx.Write("element", job.element.Render())
		}
	}
}
