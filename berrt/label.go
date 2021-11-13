package berrt

import "github.com/charmbracelet/lipgloss"

var labelmap = map[string]func(string) string{
	" ERROR ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#FF0000"),
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

	"TRACER ": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#626262"),
	).Background(
		lipgloss.Color("#262626"),
	).Bold(true).Padding(0, 1).Render,

	"RUNTIME": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Background(
		lipgloss.Color("#5F00FF"),
	).Bold(true).Padding(0, 1).Render,
}

/*
Label ...
*/
type Label struct {
	text  string
	style func(string) string
}

/*
NewLabel ...
*/
func NewLabel(text string) *Label {
	return &Label{
		text:  text,
		style: labelmap[text],
	}
}

/*
ToString returns the fully styled label as a string we can use.
*/
func (label *Label) ToString() string {
	return label.style(label.text)
}
