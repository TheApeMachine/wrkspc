package conquer

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Shell is a Platform for running commands in the underlying shell.
*/
type Shell struct {
	command []string
}

/*
Boot the runtime environment for this Platform.
*/
func (platform Shell) Boot() Platform {
	errnie.Traces()
	return platform
}

/*
Parse the command into executable steps.
*/
func (platform Shell) Parse(command []string) Platform {
	errnie.Traces()
	platform.command = command
	return platform
}

/*
Process the Command.
*/
func (platform Shell) Process() chan *spdg.Datagram {
	errnie.Traces()
	out := make(chan *spdg.Datagram)

	go func() {
		defer close(out)
		errnie.Logs("running command", platform.command[0]).With(errnie.DEBUG)
		chunks := strings.Split(platform.command[0], " ")
		platform.stream(
			// We have to properly unroll everything, otherwise we get `file not found errors`.
			exec.Command(chunks[0], chunks[1:]...),
		)

		out <- spdg.NullDatagram() // Nothing but net.
	}()

	return out
}

/*
stream executes the shell command and returns an output stream from
stdout so we can get feedback in real-time, which is needed especially
for commands that potentially never end (log streams for instance).
*/
func (platform Shell) stream(cmd *exec.Cmd) {
	errnie.Traces()

	if cmd == nil {
		return
	}

	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			errnie.Logs(scanner.Text()).With(errnie.DEBUG)
		}

		done <- struct{}{}
	}()

	errnie.Handles(cmd.Start()).With(errnie.KILL)
	<-done
	errnie.Handles(cmd.Wait()).With(errnie.KILL)
}
