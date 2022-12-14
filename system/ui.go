package system

import (
	"io"
	"os"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/tui"
)

type UIBooter struct {
	err chan error
}

func (booter *UIBooter) Kick() chan error {
	errnie.Trace()
	booter.err = make(chan error)

	dg := spd.New(spd.APPBIN, spd.UI, spd.LAYER)
	dg.Write(tui.LOGO)

	io.Copy(os.Stdout, tui.NewUI(dg))
	return booter.err
}
