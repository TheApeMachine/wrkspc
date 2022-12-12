package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/theapemachine/wrkspc/errnie"
)

type Screen struct {
	layers  []*Layer
	builder strings.Builder
	out     string
}

func NewScreen(layers ...*Layer) *Screen {
	return &Screen{layers: layers}
}

func (screen *Screen) Init() tea.Cmd {
	for _, layer := range screen.layers {
		layer.Init()
	}

	return nil
}

func (screen *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	for _, layer := range screen.layers {
		layer.Update(msg)
	}

	return screen, nil
}

func (screen *Screen) View() string {
	for _, layer := range screen.layers {
		screen.builder.WriteString(layer.View())
	}

	screen.out = screen.builder.String()
	screen.builder.Reset()
	return screen.out
}

/*
Read implements the io.Reader interface.
*/
func (screen *Screen) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented...")
	return
}

/*
Write implements the io.Writer interface.
*/
func (screen *Screen) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented...")
	return
}

/*
Close implements the io.Closer interface.
*/
func (screen *Screen) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented...")
	return errnie.Handles(nil)
}
