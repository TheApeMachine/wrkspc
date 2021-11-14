package conquer

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Command is an object that takes raw input from the command-line invocation of the program
and performs an initial aggregation of objects that will be involved in handling it.
*/
type Command struct {
	scope []string
	kube  bool
	pre   []string
	post  []string
}

/*
NewCommand constructs the wrapped command-line data into an object we can call methods on.
*/
func NewCommand(scope []string, kube bool) *Command {
	errnie.Traces()

	return &Command{
		scope: scope,
		kube:  kube,
		// The pre and post steps between which the actual command sits can be used
		// to configure the local environment and are defined in `~/.wrkspc.yml`.
		pre:  strings.Split(viper.GetString("wrkspc.run.pre"), "\n"),
		post: strings.Split(viper.GetString("wrkspc.run.post"), "\n"),
	}
}

/*
Execute the Command end-to-end.
*/
func (command *Command) Execute() chan error {
	errnie.Traces()
	out := make(chan error)

	go func() {
		defer close(out)

		// Runs a shell script from the `~/.wrkspc.yml` configuration.
		errnie.Logs("running pre command steps").With(errnie.INFO)
		command.setupAndDestroy(command.pre)

		var platform Platform

		// Select the Platform to run on which will also call Boot on that Platform so it will
		// be up and running as fast as possible.
		if !command.kube {
			platform = NewPlatform(Docker{})
		} else {
			platform = NewPlatform(Kubernetes{})
		}

		platform.Parse(command.scope).Process()

		// Runs a shell script from the `~/.wrkspc.yml` configuration.
		errnie.Logs("running post command steps").With(errnie.INFO)
		command.setupAndDestroy(command.post)

		// Send nil over the error channel to signal the program to stop.
		out <- nil
	}()

	return out
}

/*
setupAndDestroy handles the pre and post steps.
*/
func (command *Command) setupAndDestroy(stage []string) {
	for _, line := range stage {
		// No use executing an empty line, or one that is commented out.
		if string(line) == "" || string(line[0]) == "#" {
			continue
		}

		command.stream(exec.Command(line))
	}
}

/*
StreamCmd executes the shell command and returns an output stream from
stdout so we can get feedback in real-time, which is needed especially
for commands that potentially never end (log streams for instance).
*/
func (command *Command) stream(cmd *exec.Cmd) {
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
			fmt.Print(scanner.Text())
		}

		done <- struct{}{}
	}()

	errnie.Handles(cmd.Start()).With(errnie.KILL)
	<-done
	errnie.Handles(cmd.Wait()).With(errnie.KILL)
}
