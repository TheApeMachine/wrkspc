package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/system"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
runCmd is a proxy for running any terminal command using a container
which is dynamically built from an image in a configured registry.
*/
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service with the ~/.wrkspc.yml config values.",
	Long:  runtxt,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		ctx := twoface.NewContext()

		errnie.Informs("booting wrkspc for", os.Getenv("USER"))
		for err := range system.Boot(
			&system.SystemBooter{Ctx: ctx},
			//&system.KraftBooter{Ctx: ctx},
			&system.UIBooter{Ctx: ctx},
			&system.WorkspaceBooter{Ctx: ctx},
		) {
			errnie.Handles(err)
		}

		defer errnie.Ctx().Close()
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
