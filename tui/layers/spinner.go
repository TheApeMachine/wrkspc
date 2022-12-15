package layers

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	bubble   spinner.Model
	quitting bool
	err      error
}

func NewSpinner() *Spinner {
	return &Spinner{
		bubble: spinner.New(),
	}
}

func (model *Spinner) Init() tea.Cmd {
	model.bubble.Spinner = spinner.Dot
	model.bubble.Style = lipgloss.NewStyle().Foreground(
		lipgloss.Color("205"),
	)

	return model.bubble.Tick
}

func (model *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			model.quitting = true
			return model, tea.Quit
		default:
			return model, nil
		}
	case error:
		model.err = msg
		return model, nil
	default:
		var cmd tea.Cmd
		model.bubble, cmd = model.bubble.Update(msg)
		return model, cmd
	}
}

func (model *Spinner) View() string {
	if model.err != nil {
		return model.err.Error()
	}

	str := fmt.Sprintf(
		"\n\n   %s Loading forever...press q to quit\n\n",
		model.bubble.View(),
	)

	if model.quitting {
		return str + "\n"
	}

	return str
}
