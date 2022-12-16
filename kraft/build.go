package kraft

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"

	"kraftkit.sh/config"
	"kraftkit.sh/exec"
	"kraftkit.sh/pack"
	"kraftkit.sh/unikraft"

	"kraftkit.sh/iostreams"
	"kraftkit.sh/log"
	"kraftkit.sh/make"
	"kraftkit.sh/packmanager"
	"kraftkit.sh/tui/paraprogress"
	"kraftkit.sh/unikraft/app"
	"kraftkit.sh/unikraft/component"

	"kraftkit.sh/tui/processtree"
)

type BuildStage struct {
	ctx *twoface.Context
}

func (stage *BuildStage) Make() error {
	errnie.Trace()
	errnie.Informs("building unikernel with kraft")

	var err error
	var workdir string

	ctx := stage.ctx.Root()
	opts := NewOpts()
	opts.Project()

	parallel := !config.G(ctx).NoParallel
	norender := log.LoggerTypeFromString(
		config.G(ctx).Log.Type,
	) != log.FANCY

	var missingPacks []pack.Package
	var processes []*paraprogress.Process
	var searches []*processtree.ProcessTreeItem

	_, err = opts.AppCfg().Components()
	if err != nil && opts.AppCfg().Template().Name() != "" {
		packages := opts.Packages()

		search := processtree.NewProcessTreeItem(fmt.Sprintf(
			"finding %s/%s:%s...",
			opts.AppCfg().Template().Type(),
			opts.AppCfg().Template().Name(),
			opts.AppCfg().Template().Version(),
		), "", nil)

		var treemodel *processtree.ProcessTree
		if treemodel, err = processtree.NewProcessTree(
			ctx,
			[]processtree.ProcessTreeOption{
				processtree.IsParallel(parallel),
				processtree.WithRenderer(norender),
				processtree.WithFailFast(true),
			},
			search,
		); err != nil {
			return err
		}

		if err := treemodel.Start(); err != nil {
			return fmt.Errorf("could not complete search: %v", err)
		}

		proc := paraprogress.NewProcess(
			fmt.Sprintf("pulling %s", packages[0].Options().TypeNameVersion()),
			func(ctx context.Context, w func(progress float64)) error {
				return packages[0].Pull(
					ctx,
					pack.WithPullProgressFunc(w),
					pack.WithPullWorkdir(workdir),
					// pack.WithPullChecksum(!opts.NoChecksum),
					pack.WithPullCache(true),
				)
			},
		)

		processes = append(processes, proc)

		var paramodel *paraprogress.ParaProgress
		if paramodel, err = paraprogress.NewParaProgress(
			ctx,
			processes,
			paraprogress.IsParallel(parallel),
			paraprogress.WithRenderer(norender),
			paraprogress.WithFailFast(true),
		); err != nil {
			return err
		}

		if err := paramodel.Start(); err != nil {
			return fmt.Errorf("could not pull all components: %v", err)
		}
	}

	opts.Template()

	// Overwrite template with user options
	var components []component.Component
	if components, err = opts.AppCfg().Components(); err != nil {
		return errnie.Handles(err)
	}

	for _, component := range components {
		component := component // loop closure

		searches = append(searches, processtree.NewProcessTreeItem(
			fmt.Sprintf("finding %s/%s:%s...", component.Type(), component.Component().Name, component.Component().Version), "",
			func(ctx context.Context) error {
				var p []pack.Package
				if p, err = packmanager.G(ctx).Catalog(
					ctx, packmanager.CatalogQuery{
						Name: component.Name(),
						Types: []unikraft.ComponentType{
							component.Type(),
						},
						Version: component.Version(),
						NoCache: true,
					},
				); err != nil {
					return err
				}

				if len(p) == 0 {
					return fmt.Errorf("could not find: %s", component.Component().Name)
				} else if len(p) > 1 {
					return fmt.Errorf("too many options for %s", component.Component().Name)
				}

				missingPacks = append(missingPacks, p...)
				return nil
			},
		))
	}

	if len(searches) > 0 {
		treemodel, err := processtree.NewProcessTree(
			ctx,
			[]processtree.ProcessTreeOption{
				processtree.IsParallel(parallel),
				processtree.WithRenderer(norender),
				processtree.WithFailFast(true),
			},
			searches...,
		)
		if err != nil {
			return err
		}

		if err := treemodel.Start(); err != nil {
			return fmt.Errorf("could not complete search: %v", err)
		}
	}

	if len(missingPacks) > 0 {
		for _, p := range missingPacks {
			if p.Options() == nil {
				return fmt.Errorf("unexpected error occurred please try again")
			}
			p := p // loop closure
			processes = append(processes, paraprogress.NewProcess(
				fmt.Sprintf("pulling %s", p.Options().TypeNameVersion()),
				func(ctx context.Context, w func(progress float64)) error {
					return p.Pull(
						ctx,
						pack.WithPullProgressFunc(w),
						pack.WithPullWorkdir(workdir),
						pack.WithPullCache(false),
					)
				},
			))
		}

		paramodel, err := paraprogress.NewParaProgress(
			ctx,
			processes,
			paraprogress.IsParallel(parallel),
			paraprogress.WithRenderer(norender),
			paraprogress.WithFailFast(true),
		)
		if err != nil {
			return err
		}

		if err := paramodel.Start(); err != nil {
			return fmt.Errorf("could not pull all components: %v", err)
		}
	}

	processes = []*paraprogress.Process{} // reset

	targets, err := opts.AppCfg().Targets()
	if err != nil {
		return err
	}

	var mopts []make.MakeOption
	mopts = append(mopts, make.WithJobs(6))

	for _, targ := range targets {
		// See: https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		targ := targ
		processes = append(processes, paraprogress.NewProcess(
			fmt.Sprintf("configuring %s (%s)", targ.Name(), targ.ArchPlatString()),
			func(ctx context.Context, w func(progress float64)) error {
				return opts.AppCfg().DefConfig(
					ctx,
					&targ, // Target-specific options
					nil,   // No extra configuration options
					make.WithProgressFunc(w),
					make.WithSilent(true),
					make.WithExecOptions(
						exec.WithStdin(iostreams.G(ctx).In),
						exec.WithStdout(log.G(ctx).Writer()),
						exec.WithStderr(log.G(ctx).WriterLevel(logrus.ErrorLevel)),
					),
				)
			},
		))

		processes = append(processes, paraprogress.NewProcess(
			fmt.Sprintf("preparing %s (%s)", targ.Name(), targ.ArchPlatString()),
			func(ctx context.Context, w func(progress float64)) error {
				return opts.AppCfg().Prepare(
					ctx,
					&targ, // Target-specific options
					append(
						mopts,
						make.WithProgressFunc(w),
						make.WithExecOptions(
							exec.WithStdout(log.G(ctx).Writer()),
							exec.WithStderr(log.G(ctx).WriterLevel(logrus.ErrorLevel)),
						),
					)...,
				)
			},
		))

		processes = append(processes, paraprogress.NewProcess(
			fmt.Sprintf("building %s (%s)", targ.Name(), targ.ArchPlatString()),
			func(ctx context.Context, w func(progress float64)) error {
				return opts.AppCfg().Build(
					ctx,
					&targ, // Target-specific options
					app.WithBuildProgressFunc(w),
					app.WithBuildMakeOptions(append(mopts,
						make.WithExecOptions(
							exec.WithStdout(log.G(ctx).Writer()),
							exec.WithStderr(log.G(ctx).WriterLevel(logrus.ErrorLevel)),
						),
					)...),
					app.WithBuildLogFile("buildlog"),
				)
			},
		))
	}

	paramodel, err := paraprogress.NewParaProgress(
		ctx,
		processes,
		// Disable parallelization as:
		//  - The first process may be pulling the container image, which is
		//    necessary for the subsequent build steps;
		//  - The Unikraft build system can re-use compiled files from previous
		//    compilations (if the architecture does not change).
		paraprogress.IsParallel(false),
		paraprogress.WithRenderer(norender),
		paraprogress.WithFailFast(true),
	)
	if err != nil {
		return err
	}

	return paramodel.Start()
}
