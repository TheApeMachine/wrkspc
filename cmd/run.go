package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/conquer"
	"github.com/theapemachine/wrkspc/errnie"
)

var kube bool

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.PersistentFlags().BoolVarP(
		&kube, "kube", "k", false, "Run in Kubernetes cluster if true (will create one if none exists).",
	)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Proxies a command through wrkspc so it will download and run the relevant container.",
	Long:  longruntxt,
	RunE: func(cmd *cobra.Command, args []string) error {
		errnie.Traces() // If trace is true in `~/.wrkspc` output current file,
		// function, and line.
		go errnie.Runtime(30) // Same setting, print the number of goroutines every 30 secs.

		// Pass the command off to a specialist object, call Execute to set things in motion which
		// returns a channel or `error` so that can block the main goroutine and respond to
		// the error.
		return <-conquer.NewCommand(args, kube).Execute()
	},
}

var longruntxt = `
Download the relavant container that provides the command just executed if it is not already present
on the local machine, and run the container.
`
