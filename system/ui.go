package system

import (
	"io"
	"os"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/tui"
	"github.com/theapemachine/wrkspc/twoface"
)

type UIBooter struct {
	Ctx *twoface.Context
	err error
}

func (booter *UIBooter) Kick() chan error {
	errnie.Trace()
	out := make(chan error)

	go func() {
		defer close(out)
		errnie.Informs("booting ui...")

		dg := spd.New(spd.APPBIN, spd.UI, spd.LAYER)
		dg.Write(tui.LOGO)

		if _, booter.err = io.Copy(
			os.Stdout, tui.NewUI(dg),
		); errnie.Handles(booter.err) != nil {
			return
		}

		out <- nil
	}()

	return out
}
