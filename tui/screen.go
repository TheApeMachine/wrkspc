package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	term "github.com/inancgumus/screen"
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
	term.Clear()

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
Write implements the io.Writer interface.
*/
func (screen *Screen) Write(p []byte) (n int, err error) {
	return
}

/*
Close implements the io.Closer interface.
*/
func (screen *Screen) Close() error {
	return nil
}
