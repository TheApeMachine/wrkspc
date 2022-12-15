package tui

import (
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

type UI struct {
	dg      *spd.Datagram
	screens []*Screen
	builder strings.Builder
	out     string
}

func NewUI(dg *spd.Datagram) *UI {
	errnie.Trace()

	layer := dg.Next()

	ui := &UI{dg: dg, screens: []*Screen{
		NewScreen(
			NewLayer(core[string(layer)]),
		),
	}}

	if _, err := tea.NewProgram(ui).Run(); err != nil {
		errnie.Handles(err)
	}

	return ui
}

func (ui *UI) Init() tea.Cmd {
	errnie.Quiet()

	for _, screen := range ui.screens {
		screen.Init()
	}

	return nil
}

func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	errnie.Trace()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return ui, tea.Quit
		default:
			for _, screen := range ui.screens {
				screen.Update(msg)
			}
		}
	}

	return ui, nil
}

func (ui *UI) View() string {
	errnie.Trace()

	for _, screen := range ui.screens {
		ui.builder.WriteString(screen.View())
	}

	ui.out = ui.builder.String()
	ui.builder.Reset()
	return ui.out
}

func (ui *UI) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Warns("not implemented...")
	return len(p), io.EOF
}

func (ui *UI) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Warns("not implemented...")
	return
}

func (ui *UI) Close() error {
	errnie.Trace()
	errnie.Warns("not implemented...")
	return errnie.Handles(nil)
}
