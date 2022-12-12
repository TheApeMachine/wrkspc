package cmd

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/container"
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/ford"
	"github.com/theapemachine/wrkspc/gadget"
	"github.com/theapemachine/wrkspc/infra"
	"github.com/theapemachine/wrkspc/spd"
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
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(tweaker.GetBool("errnie.trace"))
		errnie.Debugging(tweaker.GetBool("errnie.debug"))
		errnie.Breakpoints(tweaker.GetBool("errnie.break"))

		gadget.NewPyroscope(
			tweaker.GetString("metrics.pyroscope.endpoint"),
		).Start()

		uipipe := drknow.NewAbstract(spd.New(
			spd.APPBIN, spd.PIPE, spd.UI,
		))

		errpipe := drknow.NewAbstract(spd.New(
			spd.APPBIN, spd.PIPE, spd.ERRNIE,
		))

		go func() {
			scanner := bufio.NewScanner(errnie.Ctx())
			for scanner.Scan() {
				os.Stdout.Write(append(scanner.Bytes(), byte('\n')))
			}
		}()

		wanpipe := drknow.NewAbstract(spd.New(
			spd.APPBIN, spd.PIPE, spd.WAN,
		))

		wrkspc := ford.NewWorkspace(
			ford.NewWorkload(
				ford.NewAssembly(errpipe),
			),
			ford.NewWorkload(
				ford.NewAssembly(uipipe),
			),
			ford.NewWorkload(
				ford.NewAssembly(wanpipe),
			),
		)

		// Only messages with the ScopeType of UI should be read from
		// a Workspace.
		io.Copy(os.Stdout, wrkspc)
		os.Exit(0)

		ctx := twoface.NewContext()
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
