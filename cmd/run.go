package cmd

import (
	"bytes"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/container"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/infra"
	"github.com/theapemachine/wrkspc/system"
	"github.com/theapemachine/wrkspc/tweaker"
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
		system.Boot(
			&system.KraftBooter{Ctx: ctx},
			&system.SystemBooter{},
			&system.UIBooter{},
			&system.WorkspaceBooter{},
		)

		os.Exit(0)

		/* TODO: Clean up... */
		brazil.NewFile(
			"~/tmp/wrkspc/quay.io/coreos/butane", "butane.yml",
			bytes.NewBuffer([]byte{}),
		)

		container.NewDocker(
			ctx, "quay.io/coreos", "butane", "latest",
		).Pull().Create(
			nil, &[]string{"butane.yml", ">", "ignition.json"},
		).Start()

		container.NewDocker(
			ctx, "pixiecore", "pixiecore", "",
		).Pull().Create(
			nil, &[]string{"boot", strings.Join([]string{
				tweaker.GetString("pxe.channel"),
				tweaker.GetString("pxe.kernel"),
			}, "/"), strings.Join([]string{
				tweaker.GetString("pxe.channel"),
				tweaker.GetString("pxe.image"),
			}, "/")},
		).Start()

		ip := tweaker.GetString("pxe.machines.deepthought.ip")
		user := tweaker.GetString("pxe.machines.deepthought.username")
		pass := tweaker.GetString("pxe.machines.deepthought.password")

		dt := infra.NewIDRAC(ip, user, pass)
		dt.Reboot()

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
