package conquer

import (
	"bytes"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/matrix"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Docker is a Platform for a container to run on without using Kubernetes.
*/
type Docker struct {
	command []string
}

/*
Boot the runtime environment for this Platform.
*/
func (platform Docker) Boot() Platform {
	errnie.Traces()
	matrix.NewDaemon() // The ContainerD daemon.
	return platform
}

/*
Parse the command into executable steps.
*/
func (platform Docker) Parse(command []string) Platform {
	errnie.Traces()
	platform.command = command
	return platform
}

/*
Process the Command.
*/
func (platform Docker) Process() chan *spdg.Datagram {
	errnie.Traces()
	out := make(chan *spdg.Datagram)

	go func() {
		disposer := twoface.NewDisposer()
		defer close(out)
		defer disposer.Cleanup()

		run := matrix.NewRun(platform.command[0], disposer)

		out <- spdg.QuickDatagram( // Send out the error wrapped into a Datagram.
			spdg.ERROR, "error", bytes.NewBuffer([]byte(
				//build.Atomic(true).Error(), // Build the image atomically and return any errors.
				errnie.Handles(run.Cycle()).ERR.Encode().Bytes(),
			)),
		)
	}()

	return out
}
