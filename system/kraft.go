package system

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/kraft"
	"github.com/theapemachine/wrkspc/twoface"
)

type KraftBooter struct {
	Ctx *twoface.Context
	err error
}

func (booter *KraftBooter) Kick() chan error {
	errnie.Trace()
	out := make(chan error)

	go func() {
		defer close(out)

		for _, stage := range kraft.NewStager(booter.Ctx) {
			if booter.err = stage.Make(); booter.err != nil {
				out <- errnie.Handles(booter.err)
			}
		}
	}()

	return out
}
