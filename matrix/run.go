package matrix

import (
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/cmd/ctr/commands"
	"github.com/containerd/containerd/cmd/ctr/commands/tasks"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Run is a wrapper that takes a container and defines a way to run it.
*/
type Run struct {
	command  string
	disposer *twoface.Disposer
}

/*
NewRun constructs an instance of Run and returns it.
TODO: Context should come from here.
*/
func NewRun(command string, disposer *twoface.Disposer) Run {
	errnie.Traces()

	return Run{
		command:  command,
		disposer: disposer,
	}
}

/*
Cycle executes a Run.
*/
func (run Run) Cycle() error {
	errnie.Traces()

	client := NewClient(Containerd{Disposer: run.disposer})
	defer client.Cleanup()

	container, spec := client.Fetch(run.command, "v1.0")
	defer container.Delete(run.disposer.Ctx, containerd.WithSnapshotCleanup)

	var (
		con console.Console
		tty = spec.Process.Terminal
	)

	if tty {
		con = console.Current()
		defer con.Reset()

		errnie.Handles(con.SetRaw()).With(errnie.KILL)
	}

	task, err := container.NewTask(
		run.disposer.Ctx, cio.NewCreator(cio.WithStdio, cio.WithTerminal),
	)
	errnie.Handles(err).With(errnie.NOOP)

	defer task.Delete(run.disposer.Ctx)
	exitStatusC, err := task.Wait(run.disposer.Ctx)
	errnie.Handles(err).With(errnie.KILL)

	// This is where we actually start the container, wrapper in an errnie Handler for a single
	// line format to be possible :)
	errnie.Handles(task.Start(run.disposer.Ctx)).With(errnie.KILL)

	if tty {
		errnie.Logs("I am consolio, I need TTY for my bunghole").With(errnie.INFO)
		errnie.Handles(tasks.HandleConsoleResize(run.disposer.Ctx, task, con))
	} else {
		sigc := commands.ForwardAllSignals(run.disposer.Ctx, task)
		defer commands.StopCatch(sigc)
	}

	errnie.Handles(task.Kill(run.disposer.Ctx, syscall.SIGTERM))

	status := <-exitStatusC

	_, _, err = status.Result()
	return errnie.Handles(err).With(errnie.KILL).ERR.First()
}
