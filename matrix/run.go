package matrix

import (
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/cmd/ctr/commands"
	"github.com/containerd/containerd/cmd/ctr/commands/tasks"
	"github.com/containerd/containerd/oci"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Run is a wrapper that takes a container and defines a way to run it.
*/
type Run struct {
	build *Build
	spec  *oci.Spec
}

/*
NewRun constructs an instance of Run and returns it.
TODO: Context should come from here.
*/
func NewRun(build *Build, spec *oci.Spec) Run {
	errnie.Traces()

	return Run{
		build: build,
		spec:  spec,
	}
}

/*
Cycle executes a Run.
*/
func (run Run) Cycle() {
	var (
		con console.Console
		tty = run.spec.Process.Terminal
	)

	if tty {
		con = console.Current()
		defer con.Reset()

		errnie.Handles(con.SetRaw()).With(errnie.KILL)
	}

	task, err := run.build.container.NewTask(run.build.disposer.Ctx, cio.NewCreator(cio.WithStdio))
	errnie.Handles(err).With(errnie.KILL)

	defer task.Delete(run.build.disposer.Ctx)
	exitStatusC, err := task.Wait(run.build.disposer.Ctx)
	errnie.Handles(err).With(errnie.KILL)

	// This is where we actually start the container, wrapper in an errnie Handler for a single
	// line format to be possible :)
	errnie.Handles(task.Start(run.build.disposer.Ctx)).With(errnie.KILL)

	if tty {
		errnie.Handles(tasks.HandleConsoleResize(run.build.disposer.Ctx, task, con))
	} else {
		sigc := commands.ForwardAllSignals(run.build.disposer.Ctx, task)
		defer commands.StopCatch(sigc)
	}

	errnie.Handles(task.Kill(run.build.disposer.Ctx, syscall.SIGTERM))

	status := <-exitStatusC

	code, _, err := status.Result()
	errnie.Handles(err).With(errnie.KILL)

	errnie.Logs(code).With(errnie.INFO)
}
