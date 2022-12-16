package kraft

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
	"kraftkit.sh/config"
	"kraftkit.sh/exec"
	"kraftkit.sh/log"
	"kraftkit.sh/machine"
	machinedriver "kraftkit.sh/machine/driver"
	machinedriveropts "kraftkit.sh/machine/driveropts"
	"kraftkit.sh/unikraft/app"
	"kraftkit.sh/unikraft/target"
)

type RunStage struct {
	ctx *twoface.Context
}

func (stage *RunStage) Make() error {
	ctx := stage.ctx.Root()
	driverType := machinedriver.QemuDriver

	var (
		store *machine.MachineStore
		err   error
	)

	store, err = machine.NewMachineStoreFromPath(
		config.G(ctx).RuntimeDir,
	)

	var driver machinedriver.Driver
	if driver, err = machinedriver.New(driverType,
		machinedriveropts.WithBackground(false),
		machinedriveropts.WithRuntimeDir(config.G(ctx).RuntimeDir),
		machinedriveropts.WithMachineStore(store),
		machinedriveropts.WithDebug(true),
		machinedriveropts.WithExecOptions(
			exec.WithStdout(os.Stdout),
			exec.WithStderr(os.Stderr),
		),
	); err != nil {
		return errnie.Handles(err)
	}

	mopts := []machine.MachineOption{
		machine.WithDriverName(driverType.String()),
		machine.WithDestroyOnExit(true),
	}

	var opts *app.ProjectOptions
	if opts, err = app.NewProjectOptions(
		nil,
		app.WithName("wrkspc"),
		app.WithWorkingDirectory(brazil.NewPath(".").Location),
		app.WithDefaultConfigPath(),
		app.WithResolvedPaths(true),
		app.WithDotConfig(false),
	); err != nil {
		return errnie.Handles(err)
	}

	var project *app.ApplicationConfig
	if project, err = app.NewApplicationFromOptions(
		opts,
	); err != nil {
		return errnie.Handles(err)
	}

	var t *target.TargetConfig
	if t, err = project.TargetByName(
		project.TargetNames()[0],
	); err != nil {
		return errnie.Handles(err)
	}

	mopts = append(mopts,
		machine.WithArchitecture(t.Architecture.Name()),
		machine.WithPlatform(t.Platform.Name()),
		machine.WithName(machine.MachineName(t.Name())),
		machine.WithAcceleration(true),
		machine.WithSource("project://"+project.Name()+":"+t.Name()),
	)

	var mid machine.MachineID
	if mid, err = driver.Create(ctx, mopts...); err != nil {
		return errnie.Handles(err)
	}

	if err = driver.Start(ctx, mid); err != nil {
		return errnie.Handles(err)
	}

	ctx, cancel := context.WithCancel(ctx)
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
	return nil
}
