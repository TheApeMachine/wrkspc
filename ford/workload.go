package ford

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

/*
A workload groups together and governs the execution of Assembly
instances.

It is responsible for facilitating communication between assemblies
whenever that is required.
*/
type Workload struct {
	assemblies []*Assembly
}

func NewWorkload(assemblies ...*Assembly) *Workload {
	errnie.Trace()

	return &Workload{
		assemblies: assemblies,
	}
}

func (wrkld *Workload) AddWork(asm ...*Assembly) *Workload {
	wrkld.assemblies = append(wrkld.assemblies, asm...)
	return wrkld
}

func (wrkld *Workload) Read(p []byte) (n int, err error) {
	errnie.Trace()

	for _, asm := range wrkld.assemblies {
		if n, err = asm.Read(p); err != nil {
			return n, errnie.Handles(err)
		}

		dg := &spd.Empty
		if err = dg.Decode(p); errnie.Handles(err) != nil {
			return
		}

		if err = wrkld.interpret(dg); errnie.Handles(err) != nil {
			return
		}
	}

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

func (wrkld *Workload) interpret(dg *spd.Datagram) error {
	errnie.Trace()
	errnie.Debugs("not implemented...")
	return errnie.Handles(nil)
}
