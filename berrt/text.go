package berrt

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

/*
Text is a message and a set of predefined styles you can apply to it.
*/
type Text struct {
	icon           string
	message        string
	mutedStyle     func(string) string
	darkStyle      func(string) string
	normalStyle    func(string) string
	highlightStyle func(string) string
}

/*
NewText returns a new Text object and predefines its styles, so we can use it to print
nicely colored text to the terminal for increased readability.
*/
func NewText(icon, message string) *Text {
	return &Text{
		icon:    icon,
		message: message,
		mutedStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#444444"),
		).Render,

		darkStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#808080"),
		).Render,

		normalStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#A8A8A8"),
		).Render,

		highlightStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#EEEEEE"),
		).Render,
	}
}

/*
ToString converts the message part of a logline into a common string format, styled with
a combination of the predefined styles.
*/
func (text *Text) ToString() string {
	return text.mutedStyle(
		time.Now().Format("2006-01-02 15:04:05.000000"),
	) + " " + NewIcon("ok") + " " + text.highlightStyle(text.message)
}
