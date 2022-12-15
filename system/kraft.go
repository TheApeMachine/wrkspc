package system

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
	"kraftkit.sh/exec"
	"kraftkit.sh/log"
	"kraftkit.sh/machine"
	machinedriver "kraftkit.sh/machine/driver"
	machinedriveropts "kraftkit.sh/machine/driveropts"
	"kraftkit.sh/unikraft/app"
	"kraftkit.sh/unikraft/target"
)

type KraftBooter struct {
	Ctx     *twoface.Context
	opts    *app.ProjectOptions
	project *app.ApplicationConfig
	err     error
}

func (booter *KraftBooter) Kick() chan error {
	errnie.Trace()
	out := make(chan error)

	go func() {
		defer close(out)

		errnie.Informs("building unikernel with kraft")

		driverType := machinedriver.QemuDriver

		var store *machine.MachineStore
		store, booter.err = machine.NewMachineStoreFromPath(
			brazil.NewPath(".").Location,
		)

		var driver machinedriver.Driver
		if driver, booter.err = machinedriver.New(driverType,
			machinedriveropts.WithBackground(false),
			machinedriveropts.WithRuntimeDir(brazil.NewPath(".").Location),
			machinedriveropts.WithMachineStore(store),
			machinedriveropts.WithDebug(true),
			machinedriveropts.WithExecOptions(
				exec.WithStdout(os.Stdout),
				exec.WithStderr(os.Stderr),
			),
		); booter.err != nil {
			out <- errnie.Handles(booter.err)
		}

		mopts := []machine.MachineOption{
			machine.WithDriverName(driverType.String()),
			machine.WithDestroyOnExit(true),
		}

		if booter.opts, booter.err = app.NewProjectOptions(
			nil,
			app.WithName("wrkspc"),
			app.WithWorkingDirectory(brazil.NewPath(".").Location),
			app.WithDefaultConfigPath(),
			app.WithResolvedPaths(true),
			app.WithDotConfig(false),
		); booter.err != nil {
			out <- errnie.Handles(booter.err)
		}

		if booter.project, booter.err = app.NewApplicationFromOptions(
			booter.opts,
		); booter.err != nil {
			out <- errnie.Handles(booter.err)
		}

		var t *target.TargetConfig
		if t, booter.err = booter.project.TargetByName(
			booter.project.TargetNames()[0],
		); booter.err != nil {
			out <- errnie.Handles(booter.err)
		}

		mopts = append(mopts,
			machine.WithArchitecture(t.Architecture.Name()),
			machine.WithPlatform(t.Platform.Name()),
			machine.WithName(machine.MachineName(t.Name())),
			machine.WithAcceleration(true),
			machine.WithSource("project://"+booter.project.Name()+":"+t.Name()),
		)

		var mid machine.MachineID
		mid, booter.err = driver.Create(booter.Ctx.Root(), mopts...)

		if booter.err = driver.Start(booter.Ctx.Root(), mid); booter.err != nil {
			out <- errnie.Handles(booter.err)
		}

		ctx, cancel := context.WithCancel(booter.Ctx.Root())
		defer cancel()

		ctrlc := make(chan os.Signal, 1)
		signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-ctrlc // wait for Ctrl+C
			fmt.Printf("\n")

			// Remove the instance on Ctrl+C if the --rm flag is passed
			log.G(ctx).Infof("removing %s...", mid.ShortString())
			if err := driver.Destroy(ctx, mid); err != nil {
				log.G(ctx).Errorf("could not remove %s: %v", mid, err)
			}

			cancel()
		}()

		go func() {
			events, errs, err := driver.ListenStatusUpdate(ctx, mid)
			if err != nil {
				ctrlc <- syscall.SIGTERM
			}

			for {
				select {
				case status := <-events:
					switch status {
					case machine.MachineStateExited, machine.MachineStateDead:
						ctrlc <- syscall.SIGTERM
						return
					}

				case err := <-errs:
					log.G(ctx).Errorf("received event error: %v", err)
					return
				}
			}
		}()

		driver.TailWriter(ctx, mid, os.Stdout)
	}()

	return out
}
