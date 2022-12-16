package kraft

import (
	"context"
	"fmt"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
	"kraftkit.sh/pack"
	"kraftkit.sh/packmanager"
	"kraftkit.sh/unikraft"
	"kraftkit.sh/unikraft/app"
)

type Opts struct {
	project  *app.ProjectOptions
	appcfg   *app.ApplicationConfig
	packages []pack.Package
	template *app.ApplicationConfig
	err      error
	ctx      context.Context
}

func NewOpts() *Opts {
	return &Opts{}
}

func (opts *Opts) Project() *app.ProjectOptions {
	if opts.project, opts.err = app.NewProjectOptions(
		nil,
		app.WithName(tweaker.Program()),
		app.WithWorkingDirectory(brazil.NewPath(".").Location),
		app.WithDefaultConfigPath(),
		app.WithResolvedPaths(true),
		app.WithDotConfig(false),
	); opts.err != nil {
		errnie.Handles(opts.err)
		return nil
	}

	return opts.project
}

func (opts *Opts) AppCfg() *app.ApplicationConfig {
	if opts.appcfg == nil {
		if opts.appcfg, opts.err = app.NewApplicationFromOptions(
			opts.project,
		); opts.err != nil {
			errnie.Handles(opts.err)
			return nil
		}
	}

	return opts.appcfg
}

func (opts *Opts) Packages() []pack.Package {
	if opts.packages == nil || len(opts.packages) == 0 {
		if opts.packages, opts.err = packmanager.G(opts.ctx).Catalog(
			opts.ctx, packmanager.CatalogQuery{
				Name: opts.AppCfg().Template().Name(),
				Types: []unikraft.ComponentType{
					unikraft.ComponentTypeApp,
				},
				Version: opts.AppCfg().Template().Version(),
				NoCache: true,
			},
		); opts.err != nil {
			return nil
		}

		if len(opts.packages) == 0 {
			errnie.Handles(fmt.Errorf(
				"could not find: %s", opts.appcfg.Template().Name(),
			))
		}

		if len(opts.packages) > 1 {
			errnie.Handles(fmt.Errorf(
				"too many options for %s", opts.appcfg.Template().Name(),
			))
		}
	}

	return opts.packages
}

func (opts *Opts) Template() *app.ProjectOptions {
	if opts.AppCfg().Template().Name() != "" {
		var templateWorkdir string
		if templateWorkdir, opts.err = unikraft.PlaceComponent(
			brazil.NewPath(".").Location,
			opts.AppCfg().Template().Type(),
			opts.AppCfg().Template().Name(),
		); errnie.Handles(opts.err) != nil {
			return nil
		}

		var templateOps *app.ProjectOptions
		if templateOps, opts.err = app.NewProjectOptions(
			nil,
			app.WithWorkingDirectory(templateWorkdir),
			app.WithDefaultConfigPath(),
			app.WithResolvedPaths(true),
			app.WithDotConfig(false),
		); errnie.Handles(opts.err) != nil {
			return nil
		}

		if opts.template, opts.err = app.NewApplicationFromOptions(
			templateOps,
		); errnie.Handles(opts.err) != nil {
			return nil
		}

		opts.appcfg = opts.template.MergeTemplate(opts.appcfg)
	}

	return opts.project
}
