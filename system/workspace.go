package system

import (
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/ford"
	"github.com/theapemachine/wrkspc/spd"
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
		spinner := spd.New(spd.APPBIN, spd.UI, spd.SPINNER)
		spinner.Write([]byte("loading your wrkspc"))

		ford.NewWorkspace(
			ford.NewWorkload(
				ford.NewAssembly(
					drknow.NewAbstract(spinner),
				),
			),
		)
	}()

	return out
}
