package matrix

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/pkg/seed"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/urfave/cli"
)

func init() {
	seed.WithTimeAndRand()
}

/*
NewDaemon starts a new ContainerD daemon.
*/
func NewDaemon() {
	errnie.Traces()

	go func() {
		// shim.Run("io.containerd.runc.v2", v2.New)
		daemon := command.App()
		daemon.Flags = []cli.Flag{
			cli.StringFlag{
				Name:  "config,c",
				Usage: "path to the configuration file",
				Value: filepath.Join(defaults.DefaultConfigDir, "config.toml"),
			},
			cli.StringFlag{
				Name:  "log-level,l",
				Usage: "set the logging level [trace, debug, info, warn, error, fatal, panic]",
				Value: "panic",
			},
			cli.StringFlag{
				Name:  "address,a",
				Usage: "address for containerd's GRPC server",
			},
			cli.StringFlag{
				Name:  "root",
				Usage: "containerd root directory",
			},
			cli.StringFlag{
				Name:  "state",
				Usage: "containerd state directory",
			},
		}

		for _, f := range daemon.Flags {
			if f.GetName() == "log-level,l" {
				f.Apply(flag.NewFlagSet("-l error", flag.ExitOnError))
			}
		}

		errnie.Handles(
			daemon.Run(os.Args),
		).With(errnie.KILL)
	}()
}
