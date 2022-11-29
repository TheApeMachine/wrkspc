package cmd

import (
	"context"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/spd"
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
		server := spd.AnnounceServer{}
		client := spd.Announce_ServerToClient(server)

		conn := rpc.NewConn(rpc.NewStreamTransport(
			spd.New([]byte("test"), []byte("test"), []byte("test"), nil),
		), &rpc.Options{
			// The BootstrapClient is the RPC interface that will be made available
			// to the remote endpoint by default.  In this case, Arith.
			BootstrapClient: capnp.Client(client),
		})

		ctx := context.Background()
		ctx, cnl := context.WithCancel(ctx)

		defer conn.Close()
		select {
		case <-conn.Done():
			cnl()
			return nil
		case <-ctx.Done():
			return conn.Close()
		}
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
