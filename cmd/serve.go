package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/bcknd"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve wrkspc as a service.",
	Long:  longservetxt,
	RunE: func(cmd *cobra.Command, args []string) error {
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
