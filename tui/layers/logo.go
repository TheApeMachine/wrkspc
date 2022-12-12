package layers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type Logo struct {
	out string
}

func NewLogo() *Logo {
	return &Logo{}
}

func (logo *Logo) Init() tea.Cmd {
	logo.out = lipgloss.Place(
		96, 9, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(
			figure.NewColorFigure(
				"WRKSPC", "isometric1", "purple", true,
			).String(),
		),
		lipgloss.WithWhitespaceChars("猫咪"),
		lipgloss.WithWhitespaceForeground(
			lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"},
		),
	)

	return nil
}

func (logo *Logo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return logo, nil
}

func (logo *Logo) View() string {
	return logo.out
}
