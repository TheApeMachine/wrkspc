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
		&kube, "kube", "k", true, "Run in Kubernetes cluster if true (will create one if none exists).",
	)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Proxies a command through wrkspc so it will download and run the relevant container.",
	Long:  longruntxt,
	RunE: func(cmd *cobra.Command, args []string) error {
		errnie.Traces()
		go errnie.Runtime(30)
		command := conquer.NewCommand(args[0], kube)
		return <-command.Execute()
	},
}

var longruntxt = `
Download the relavant container that provides the command just executed if it is not already present
on the local machine, and run the container.
`
