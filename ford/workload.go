package ford

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
A workload groups together and governs the execution of Assembly
instances.

It is responsible for facilitating communication between assemblies
whenever that is required.
*/
type Workload struct {
	io.ReadWriteCloser
	assemblies []*Assembly
}

func NewWorkload(assemblies []*Assembly) *Workload {
	errnie.Trace()
	return &Workload{
		assemblies: assemblies,
	}
}

func (wrkld *Workload) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (wrkld *Workload) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (wrkld *Workload) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}