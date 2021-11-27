package cmd

import (
	"os"
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
		// We can leave the name input empty for now since we went with writing out the binaries
		// for now. Execution in-memory directly is technically working, but error output is not
		// great, so runc seems to be wanting some flags or something.
		binner := matroesjka.NewEmbed("")
		binner.Write() // Write out the embedded binaries.

		// We want to override te executable paths of the user for a while so we contain them to
		// only the embdded tooling, such that we affect the system in the least possible way.
		oldpath := os.Getenv("PATH")
		os.Setenv("PATH", brazil.HomePath()+"/wrkspc")
		defer os.Setenv("PATH", oldpath)

		// Start the docker daemon.
		conquer.NewCommand([]string{
			brazil.HomePath() + "/wrkspc/dockerd",
		}, conquer.SHELL).Execute()

		time.Sleep(5 * time.Second)

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
