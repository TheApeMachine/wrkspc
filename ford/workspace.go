package ford

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

type Workspace struct {
	io.ReadWriteCloser
	workloads []*Workload
	ctx       twoface.Context
	pool      *twoface.Pool
}

func NewWorkspace(workloads ...*Workload) *Workspace {
	errnie.Trace()
	ctx := twoface.NewContext(nil)

	return &Workspace{
		workloads: workloads,
		ctx:       ctx,
		pool:      twoface.NewPool(ctx).Run(),
	}
}

func (wrkspc *Workspace) Read(p []byte) (n int, err error) {
	errnie.Trace()
	return wrkspc.ctx.Read(p)
}

func (wrkspc *Workspace) Write(p []byte) (n int, err error) {
	errnie.Trace()
	return wrkspc.ctx.Write(p)
}

func (wrkspc *Workspace) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}
