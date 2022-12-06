package ford

import (
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
A workload groups together and governs the execution of Assembly
instances.

It is responsible for facilitating communication between assemblies
whenever that is required.
*/
type Assembly struct {
	abstracts []*drknow.Abstract
	size      int
}

func NewAssembly(abstracts ...*drknow.Abstract) *Assembly {
	errnie.Trace()

	return &Assembly{
		abstracts: abstracts,
		size:      len(abstracts),
	}
}

func (asm *Assembly) Read(p []byte) (n int, err error) {
	errnie.Trace()
	return
}

func (asm *Assembly) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return
}

func (asm *Assembly) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return errnie.NewError(nil)
}
