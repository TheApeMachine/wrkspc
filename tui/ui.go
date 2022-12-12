package tui

import (
	"io"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/tui/layers"
)

type UI struct {
	datagram *spd.Datagram
	screens  []*Screen
	builder  strings.Builder
	out      string
}

func NewUI(datagram *spd.Datagram) *UI {
	ui := &UI{
		datagram: datagram,
		screens:  []*Screen{NewScreen(NewLayer(layers.NewLogo()))},
	}

	if _, err := tea.NewProgram(ui).Run(); err != nil {
		log.Panic(err)
	}

	return ui
}

func (ui *UI) Init() tea.Cmd {
	for _, screen := range ui.screens {
		screen.Init()
	}

	return nil
}

func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	for _, screen := range ui.screens {
		screen.Update(msg)
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
	return len(p), io.EOF
}

func (ui *UI) Write(p []byte) (n int, err error) {
	return
}

func (ui *UI) Close() error {
	return errnie.Handles(nil)
}
