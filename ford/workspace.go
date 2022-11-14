package ford

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Workspace is a mechanism to orchestrate many objects,
using a twoface.Context as a stateful object to
communicate data around.
*/
type Workspace struct {
	ctx      twoface.Context
	pool     *twoface.Pool
	assembly *Assembly
}

/*
NewWorkspace returns a pointer to a new instance
of a Workspace.
*/
func NewWorkspace(ctx twoface.Context) *Workspace {
	return &Workspace{
		ctx:      ctx,
		pool:     twoface.NewPool(ctx).Run(),
		assembly: NewAssembly(ctx),
	}
}

func (workspace *Workspace) Add(workload Workload) errnie.Error {
	return workspace.assembly.Add(workload)
}

func (workspace *Workspace) Do() {
	for _, workload := range workspace.assembly.Workloads {
		workload.Read([]byte{})
	}
}
