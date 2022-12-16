package kraft

import (
	"context"

	"github.com/theapemachine/wrkspc/twoface"
	"kraftkit.sh/config"
	"kraftkit.sh/log"
	"kraftkit.sh/packmanager"
	"kraftkit.sh/tui/processtree"
)

type PkgStage struct {
	ctx *twoface.Context
}

func (stage *PkgStage) Make() error {
	ctx := stage.ctx.Root()
	pm := packmanager.G(ctx)

	parallel := !config.G(ctx).NoParallel
	norender := log.LoggerTypeFromString(config.G(ctx).Log.Type) != log.FANCY

	model, err := processtree.NewProcessTree(
		ctx,
		[]processtree.ProcessTreeOption{
			// processtree.WithVerb("Updating"),
			processtree.IsParallel(parallel),
			processtree.WithRenderer(norender),
		},
		[]*processtree.ProcessTreeItem{
			processtree.NewProcessTreeItem(
				"Updating...",
				"",
				func(ctx context.Context) error {
					return pm.Update(ctx)
				},
			),
		}...,
	)
	if err != nil {
		return err
	}

	if err := model.Start(); err != nil {
		return err
	}

	return nil
}
