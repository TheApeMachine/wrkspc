package system

import (
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
	"kraftkit.sh/unikraft/app"
)

type KraftBooter struct {
	Ctx     *twoface.Context
	opts    *app.ProjectOptions
	project *app.ApplicationConfig
	err     error
}

func (booter *KraftBooter) Kick() chan error {
	errnie.Informs("building unikernel with kraft")
	out := make(chan error)

	if booter.opts, booter.err = app.NewProjectOptions(
		nil,
		app.WithWorkingDirectory(brazil.NewPath(".").Location),
		app.WithDefaultConfigPath(),
		app.WithResolvedPaths(true),
		app.WithDotConfig(false),
	); booter.err != nil {
		out <- errnie.Handles(booter.err)
	}

	if booter.project, booter.err = app.NewApplicationFromOptions(
		booter.opts,
	); booter.err != nil {
		out <- errnie.Handles(booter.err)
	}

	booter.project.Configure(booter.Ctx.Root(), nil)

	return out
}
