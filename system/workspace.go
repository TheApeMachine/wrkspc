package system

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/ford"
	"github.com/theapemachine/wrkspc/twoface"
)

type WorkspaceBooter struct {
	Ctx *twoface.Context
	err error
}

func (booter *WorkspaceBooter) Kick() chan error {
	errnie.Trace()
	out := make(chan error)

	go func() {
		defer close(out)
		errnie.Informs("booting workspace...")

		ford.NewWorkspace(
			ford.NewWorkload(
				ford.NewAssembly(),
			),
		)

		out <- nil
	}()

	return out
}
