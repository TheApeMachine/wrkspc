package ford

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

type Workspace struct {
	io.ReadWriteCloser
	workloads []*Workload
	ctx       *twoface.Context
}

func NewWorkspace(workloads ...*Workload) *Workspace {
	errnie.Trace()

	return &Workspace{
		workloads: workloads,
	}
}

func (wrkspc *Workspace) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (wrkspc *Workspace) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (wrkspc *Workspace) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}
