package ford

import (
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
)

var Main *Workspace

func init() {
	Main = NewWorkspace()
}

type Workspace struct {
	ID        uuid.UUID
	tree      *drknow.Tree
	workloads []*Workload
	err       chan error
}

func NewWorkspace(workloads ...*Workload) *Workspace {
	errnie.Trace()

	return &Workspace{
		uuid.New(),
		drknow.NewTree(),
		workloads,
		make(chan error),
	}
}

func (wrkspc *Workspace) AddWork(work ...*Workload) *Workspace {
	wrkspc.workloads = append(wrkspc.workloads, work...)
	return wrkspc
}

func (wrkspc *Workspace) Read(p []byte) (n int, err error) {
	errnie.Trace()

	for _, workload := range wrkspc.workloads {
		if n, err = workload.Read(p); err != nil {
			wrkspc.Close()
		}
	}

	return
}

func (wrkspc *Workspace) Write(p []byte) (n int, err error) {
	errnie.Trace()

	for _, workload := range wrkspc.workloads {
		if n, err = workload.Write(p); err != nil {
			wrkspc.Close()
		}
	}

	return
}

func (wrkspc *Workspace) Close() error {
	errnie.Trace()
	errnie.Informs("closing workspace", wrkspc.ID)

	for _, workload := range wrkspc.workloads {
		workload.Close()
	}

	wrkspc.err <- errnie.NewError(nil)
	return nil
}
