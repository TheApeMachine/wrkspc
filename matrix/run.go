package matrix

import (
	"context"
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
	errnie.Logs("build & spec", build, spec).With(errnie.DEBUG)

	return Run{
		build: build,
		spec:  spec,
	}
}

/*
Cycle executes a Run.
*/
func (run Run) Cycle(ctx context.Context) error {
	errnie.Traces()

	var (
		con console.Console
		tty = run.spec.Process.Terminal
	)

	if tty {
		con = console.Current()
		defer con.Reset()

		errnie.Handles(con.SetRaw()).With(errnie.KILL)
	}

	task, err := run.build.container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	errnie.Handles(err).With(errnie.KILL)
	errnie.Logs("task", task).With(errnie.INFO)

	defer task.Delete(run.build.disposer.Ctx)
	exitStatusC, err := task.Wait(run.build.disposer.Ctx)
	errnie.Handles(err).With(errnie.KILL)

	// This is where we actually start the container, wrapper in an errnie Handler for a single
	// line format to be possible :)
	errnie.Handles(task.Start(run.build.disposer.Ctx)).With(errnie.KILL)

	if tty {
		errnie.Logs("I am consolio, I need TTY for my bunghole").With(errnie.INFO)
		errnie.Handles(tasks.HandleConsoleResize(run.build.disposer.Ctx, task, con))
	} else {
		errnie.Logs("running without TTY").With(errnie.INFO)
		sigc := commands.ForwardAllSignals(run.build.disposer.Ctx, task)
		defer commands.StopCatch(sigc)
	}

	errnie.Handles(task.Kill(run.build.disposer.Ctx, syscall.SIGTERM))

	status := <-exitStatusC

	code, _, err := status.Result()
	errnie.Handles(err).With(errnie.KILL)

	errnie.Logs(code).With(errnie.INFO)
	return err
}
