package ford

import (
	"io"

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
	io.ReadWriteCloser
	abstract drknow.Abstract
}

func NewAssembly(abstract drknow.Abstract) *Assembly {
	errnie.Trace()

	return &Assembly{
		abstract: abstract,
	}
}

func (asm *Assembly) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

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
