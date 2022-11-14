package ford

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Assembly is a collection of Part that will be iterated over using
a twoface.Pool where each Part is a twoface.Job.
*/
type Assembly struct {
	io.ReadWriter
	twoface.Job

	ctx       twoface.Context
	Workloads []Workload
}

/*
NewAssembly returns a pointer to a new instance of Assembly.
*/
func NewAssembly(ctx twoface.Context) *Assembly {
	return &Assembly{
		ctx:       ctx,
		Workloads: make([]Workload, 0),
	}
}

func (assembly *Assembly) Add(workload Workload) errnie.Error {
	assembly.Workloads = append(
		assembly.Workloads,
		workload,
	)

	return errnie.NewError(nil)
}
