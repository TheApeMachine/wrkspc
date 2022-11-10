package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/docker"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/infra"
	"github.com/theapemachine/wrkspc/kube"
	"github.com/theapemachine/wrkspc/twoface"
)

var orchestrator string

var orchestratorMap = map[string]func() infra.Cluster{
	// "nomad":      nomad.NewCluster,
	"kubernetes": kube.NewCluster,
}

var clientMap = map[string]func() infra.Client{
	// "nomad":      nomad.NewClient,
	"kubernetes": kube.NewClient,
}

/*
init will run before anything else (including main function).
*/
func init() {
	// Add a new command to the Cobra Commander root command
	// (which is the compiled binary).
	rootCmd.AddCommand(runCmd)
	rootCmd.PersistentFlags().StringVarP(
		&orchestrator, "orchestrator", "o", "kubernetes",
		"The orchestrator to use <nomad|kubernetes>.",
	)
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
		// Set verbosity level for errnie.
		errnie.Tracing(true)
		errnie.Debugging(true)

		// Setup os interrupt signal handling.
		signals := twoface.NewSignal()
		stop := signals.Run()

		// Get a handle on a new cluster object.
		cluster := orchestratorMap[orchestrator]()

		// Provision the cluster and report back an errnie.Error, which we
		// convert to a native Go error before returning.
		if err := cluster.Provision(); err.Type != errnie.NIL {
			// TODO: It is not really correct to assume any error here just
			// means our cluster already exists, but for now this should not
			// present any issues.
			errnie.Logs("cluster already provisioned").With(
				errnie.SUCCESS,
			)
		}

		client := clientMap[orchestrator]()

		client.Apply("system", "system", "system")
		client.Apply("base", "istio", "istio-system")
		client.Apply("istiod", "istio", "istio-system")
		// client.Apply("vault", "hashicorp", "vault")
		// client.Apply("minio", "minio", "minio")
		client.Apply("harbor", "harbor", "harbor")

		docker.NewBuilder("wrkgrp", "wrkspc", "test").Build()

		client.Apply("gateway-service", "wrkspc", "services")
		client.Apply("lumiere-service", "wrkspc", "services")

		// Send a message to the interrupt handler to stop the program.
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
