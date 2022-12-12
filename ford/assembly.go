package ford

import (
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

/*
A workload groups together and governs the execution of Assembly
instances.

It is responsible for facilitating communication between assemblies
whenever that is required.
*/
type Assembly struct {
	abstracts []*drknow.Abstract
}

func NewAssembly(abstracts ...*drknow.Abstract) *Assembly {
	errnie.Trace()
	return &Assembly{abstracts}
}

func (asm *Assembly) Read(p []byte) (n int, err error) {
	errnie.Trace()

	for _, abstract := range asm.abstracts {
		if n, err = abstract.Read(p); err != nil {
			return n, errnie.Handles(err)
		}

		dg := &spd.Empty
		if err = dg.Decode(p); errnie.Handles(err) != nil {
			return
		}

		if err = asm.interpret(dg); errnie.Handles(err) != nil {
			return
		}
	}

	return
}

func (asm *Assembly) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented...")
	return
}

func (asm *Assembly) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented...")

	return errnie.NewError(nil)
}

func (asm *Assembly) interpret(dg *spd.Datagram) error {
	errnie.Trace()
	errnie.Debugs("not implemented...")
	return errnie.Handles(nil)
}
