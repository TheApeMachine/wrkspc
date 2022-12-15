package layers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
	"github.com/theapemachine/wrkspc/errnie"
)

type Logo struct {
	out string
}

func NewLogo() *Logo {
	errnie.Trace()
	return &Logo{}
}

func (logo *Logo) Init() tea.Cmd {
	errnie.Trace()

	logo.out = figure.NewColorFigure(
		"WRKSPC", "isometric1", "purple", true,
	).String()

	return nil
}

func (logo *Logo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	errnie.Trace()

	return logo, nil
}

func (logo *Logo) View() string {
	errnie.Trace()
	return logo.out
}
