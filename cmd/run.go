package cmd

import (
	"bufio"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/ford"
	"github.com/theapemachine/wrkspc/gadget"
	"github.com/theapemachine/wrkspc/tweaker"
)

var orchestrator string

/*
runCmd is a proxy for running any terminal command using a container
which is dynamically built from an image in a configured registry.
*/
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service with the ~/.wrkspc.yml config values.",
	Long:  runtxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(tweaker.GetBool("errnie.trace"))
		errnie.Debugging(tweaker.GetBool("errnie.debug"))
		errnie.Breakpoints(tweaker.GetBool("errnie.break"))

		gadget.NewPyroscope(
			tweaker.GetString("metrics.pyroscope.endpoint"),
		).Start()

		wrkspc := bufio.NewReader(ford.NewWorkspace())
		return <-sockpuppet.NewChannel(wrkspc)
	},
}

/*
runtxt lives here to keep the command definition section cleaner.
*/
var runtxt = `
Use this sub command to proxy any terminal command through and it will
look for an existing image in the configured registry which has the command
included, build that image into a container and deploy it onto the
Kubernetes cluster that will be created first.
`
