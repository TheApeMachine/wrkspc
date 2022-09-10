package tui

import "github.com/charmbracelet/lipgloss"

var labelMap = map[string]func(string) string{
	" TRACE ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#626262"),
	).Background(
		lipgloss.Color("#262626"),
	).Bold(true).Padding(0, 1).Render,

	"RUNTIME": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#8300FF"),
	).Bold(true).Padding(0, 1).Render,

	" NOOP  ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#FF0055"),
	).Bold(true).Padding(0, 1).Render,

	" KILL  ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#FF0055"),
	).Bold(true).Padding(0, 1).Render,

	"SUCCESS": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#00FF55"),
	).Bold(true).Padding(0, 1).Render,

	" INFO  ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#0087FF"),
	).Bold(true).Padding(0, 1).Render,

	" DEBUG ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#262626"),
	).Background(
		lipgloss.Color("#626262"),
	).Bold(true).Padding(0, 1).Render,

	"WARNING": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#FFAF00"),
	).Bold(true).Padding(0, 1).Render,

	" ERROR ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#FF0055"),
	).Bold(true).Padding(0, 1).Render,
}

type Label struct {
	val string
}

func NewLabel(val string) Label {
	return Label{val: val}
}

func (label Label) Print() string {
	return labelMap[label.val](label.val)
}
