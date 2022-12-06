package ford

import (
	"github.com/theapemachine/wrkspc/errnie"
)

type Workspace struct {
	workloads []*Workload
}

func NewWorkspace(workloads ...*Workload) *Workspace {
	errnie.Trace()
	return &Workspace{workloads}
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
