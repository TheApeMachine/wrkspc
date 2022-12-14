package system

import (
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/ford"
	"github.com/theapemachine/wrkspc/spd"
)

type WorkspaceBooter struct {
	err chan error
}

func (booter *WorkspaceBooter) Kick() chan error {
	errnie.Trace()

	ford.NewWorkspace(
		ford.NewWorkload(
			ford.NewAssembly(drknow.NewAbstract(
				spd.New(spd.APPBIN, spd.SYSTEM, spd.BOOT),
			)),
		),
	)

	return booter.err
}
