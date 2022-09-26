package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/errnie"
)

var conns string
var stores string
var jobs string

/*
init will run before anything else (including main function).
*/
func init() {
	// Add a new command to the Cobra Commander root command
	// (which is the compiled binary).
	rootCmd.AddCommand(serveCmd)
	rootCmd.PersistentFlags().StringVarP(
		&conns, "conns", "c", "http",
		"The available network connections for the service.",
	)

	rootCmd.PersistentFlags().StringVarP(
		&stores, "stores", "s", "s3,radix",
		"The available stores for the service.",
	)

	rootCmd.PersistentFlags().StringVarP(
		&conns, "conns", "c", "http",
		"The available network connections for the service.",
	)
}

/*
serveCmd is a proxy for running any terminal command using a container
which is dynamically built from an image in a configured registry.
*/
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the service",
	Long:  servetxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(true)
		errnie.Debugging(true)

		return nil
	},
}

/*
runtxt lives here to keep the command definition section cleaner.
*/
var servetxt = `
TODO: explanation.
`
