package layers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
)

type Logo struct {
	out string
}

func NewLogo() *Logo {
	return &Logo{}
}

func (logo *Logo) Init() tea.Cmd {
	logo.out = figure.NewColorFigure(
		"WRKSPC", "isometric1", "purple", true,
	).String()

	return nil
}

func (logo *Logo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []byte:
		fmt.Println(string(msg))
	}

	return logo, nil
}

func (logo *Logo) View() string {
	return logo.out
}
