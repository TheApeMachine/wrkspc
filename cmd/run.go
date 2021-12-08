package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/conquer"
	"github.com/theapemachine/wrkspc/dbcooper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/matroesjka"
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
		errnie.Runtime(30) // Same setting, print the number of goroutines every 30 secs.

		// We can leave the name input empty for now since we went with writing out the binaries
		// for now. Execution in-memory directly is technically working, but error output is not
		// great, so runc seems to be wanting some flags or something.
		binner := matroesjka.NewEmbed()
		binner.Write() // Write out the embedded binaries.

		// Shield ourselves into a limited environment that we build, separate of the user's
		// underlying environment.
		hijacker := dbcooper.NewHijacker().TakeOver()
		defer hijacker.Release()

		// Pass the command off to a specialist object, call Execute to set things in motion which
		// returns a channel of `error` so that can block the main goroutine, which will unblock
		// as soon as you push an error onto the channel.
		return <-conquer.NewCommand(args, conquer.DOCKER).Execute()
	},
}

var longruntxt = `
Download the relavant container that provides the command just executed if it is not already present
on the local machine, and run the container.
`
