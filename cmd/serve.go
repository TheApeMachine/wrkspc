package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/bcknd"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/conquer"
	"github.com/theapemachine/wrkspc/matroesjka"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve wrkspc as a service.",
	Long:  longservetxt,
	RunE: func(cmd *cobra.Command, args []string) error {
		binner := matroesjka.NewEmbed("runc")
		binner.Write()

		conquer.NewCommand(
			[]string{brazil.HomePath() + "/wrkspc/containerd"}, conquer.SHELL,
		).Execute()

		time.Sleep(3 * time.Second)

		conquer.NewCommand([]string{
			brazil.HomePath() + "/wrkspc/containerd-shim-runc-v2 -namespace moby -id 1bc362c60101a8077bc7b7748f7127fc18d760f6d4b2fde1b7199cf957523476 -address /run/containerd/containerd.sock",
		}, conquer.SHELL).Execute()

		time.Sleep(3 * time.Second)

		conquer.NewCommand([]string{
			brazil.HomePath() + "/wrkspc/dockerd -H fd:// --containerd=/run/containerd/containerd.sock",
		}, conquer.SHELL).Execute()

		time.Sleep(3 * time.Second)

		conquer.NewCommand([]string{
			brazil.HomePath() + "/wrkspc/docker-proxy -proto tcp -host-ip 0.0.0.0 -host-port 5900 -container-ip 172.17.0.2 -container-port 5900",
		}, conquer.SHELL).Execute()

		time.Sleep(3 * time.Second)

		conquer.NewCommand([]string{
			brazil.HomePath() + "/wrkspc/docker-proxy -proto tcp -host-ip 0.0.0.0 -host-port 5800 -container-ip 172.17.0.2 -container-port 5800",
		}, conquer.SHELL).Execute()

		time.Sleep(3 * time.Second)

		switch args[0] {
		case "bcknd":
			for dg := range bcknd.NewServer().Up() {
				_ = dg
			}
		}

		return nil
	},
}

var longservetxt = `
The serve command can launch either the built in micro services, or take any container name and
pull that as a service from your container registry.

Built in services are:

- bcknd : Starts a server to connect to for ingress and egress.
- srch  : Searches anything that can be searched.
- mmry  : Provides an in-memory search index with some additional features.

These should be able to function as a fully featured backend system that coover most use-cases.
`
