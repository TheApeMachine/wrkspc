package brian

import "github.com/theapemachine/wrkspc/twoface"

type Wave interface {
	Gossip(*twoface.Context) *twoface.Context
}

func NewWave(waveType Wave) Wave {
	return waveType
}

type ProtoWave struct {
	ctx *twoface.Context
}

func NewProtoWave(ctx *twoface.Context) Wave {
	return ProtoWave{
		ctx: ctx,
	}
}

func (wave ProtoWave) Gossip(ctx *twoface.Context) *twoface.Context {
	return ctx
}
