package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/kube"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
init will run before anything else (including main function).
*/
func init() {
	// Add a new command to the Cobra Commander root command
	// (which is the compiled binary).
	rootCmd.AddCommand(runCmd)
}

/*
runCmd is a proxy for running any terminal command using a container
which is dynamically built from an image in a configured registry.
*/
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service with the ~/.wrkspc.yml config values.",
	Long:  runtxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(true)
		errnie.Debugging(true)

		signals := twoface.NewSignal()
		stop := signals.Run()

		// Get a handle on a new KIND (Kubernetes In Docker) cluster object.
		cluster := kube.NewCluster(kube.KIND)

		// Provision the cluster and report back an errnie.Error, which we
		// convert to a native Go error before returning.
		if err := cluster.Provision(); err.Type != errnie.NIL {
			// TODO: It is not really correct to assume any error here just
			// means our cluster already exists, but for now this should not
			// present any issues.
			errnie.Logs("cluster already provisioned").With(errnie.SUCCESS)
			cluster.IsProvisioned = true
		}

		stop <- struct{}{}

		return nil
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
