package tui

import (
	"io"

	"github.com/common-nighthawk/go-figure"
	"github.com/theapemachine/wrkspc/errnie"
)

type UI struct {
	title string
}

func NewUI(title string) *UI {
	return &UI{title}
}

func (ui *UI) Read(p []byte) (n int, err error) {
	figure.NewColorFigure(ui.title, "isometric1", "purple", true).Print()
	return len(p), io.EOF
}

func (ui *UI) Write(p []byte) (n int, err error) {
	return
}

func (ui *UI) Close() error {
	return errnie.Handles(nil)
}
