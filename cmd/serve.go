package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/passepartout"
	"github.com/theapemachine/wrkspc/zaha"
)

var service string
var port string

/*
init will run before anything else (including main function).
*/
func init() {
	// Add a new command to the Cobra Commander root command
	// (which is the compiled binary).
	rootCmd.AddCommand(serveCmd)

	rootCmd.PersistentFlags().StringVarP(
		&service, "service", "s", "gateway",
		"The service that should be constructed.",
	)

	rootCmd.PersistentFlags().StringVarP(
		&port, "port", "p", "8090",
		"The port the service should run on.",
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

		architecture := zaha.NewArchitecture(
			service,
			zaha.NewNetwork(nil, nil),
			[]passepartout.Store{},
		)
		return architecture.Do()
	},
}

/*
runtxt lives here to keep the command definition section cleaner.
*/
var servetxt = `
TODO: explanation.
`
