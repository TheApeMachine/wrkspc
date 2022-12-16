package tui

import (
	"log"
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
	layer := dg.Next()

	ui := &UI{dg: dg, screens: []*Screen{
		NewScreen(
			NewLayer(core[string(layer)]),
		),
	}}

	errnie.Quiet(ui)

	go func() {
		if _, err := tea.NewProgram(ui).Run(); err != nil {
			log.Fatalln(err)
		}
	}()

	return ui
}

func (ui *UI) Init() tea.Cmd {
	for _, screen := range ui.screens {
		screen.Init()
	}

	return nil
}

func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return ui, tea.Quit
		}
	default:
		for _, screen := range ui.screens {
			screen.Update(msg)
		}
	}

	return ui, nil
}

func (ui *UI) View() string {
	for _, screen := range ui.screens {
		ui.builder.WriteString(screen.View())
	}

	ui.out = ui.builder.String()
	ui.builder.Reset()
	return ui.out
}

func (ui *UI) Read(p []byte) (n int, err error) {
	return
}

func (ui *UI) Write(p []byte) (n int, err error) {
	ui.Update(tea.Msg(p))
	return
}

func (ui *UI) Close() error {
	return nil
}
