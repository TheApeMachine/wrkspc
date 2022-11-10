package container

import (
	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/urfave/cli"
)

/*
Daemon is a wrapper around the ContainerD daemon.
*/
type Daemon struct {
	app *cli.App
	err errnie.Error
}

/*
NewDaemon returns a pointer to a new instance of Daemon.
*/
func NewDaemon() *Daemon {
	errnie.Traces()

	return &Daemon{
		app: command.App(),
	}
}

/*
Run the Daemon instance and start a ContainerD process.
*/
func (daemon *Daemon) Run() *Daemon {
	daemon.err = errnie.Handles(daemon.app.Run([]string{"run"}))
	return daemon
}

/*
Error implements the Go error interface.
*/
func (daemon *Daemon) Error() string {
	return daemon.err.Msg
}
