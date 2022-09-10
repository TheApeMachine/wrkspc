package tui

import "github.com/charmbracelet/lipgloss"

var colorMap = map[string]func(string) string{
	"MUTE": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#444444"),
	).Align(lipgloss.Left).Render,

	"DARK": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#808080"),
	).Align(lipgloss.Left).Render,

	"NORM": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#A8A8A8"),
	).Align(lipgloss.Left).Render,

	"HIGH": lipgloss.NewStyle().Foreground(
		lipgloss.Color("#EEEEEE"),
	).Align(lipgloss.Left).Render,
}

type Color struct {
	clr string
	val string
}

func NewColor(clr, val string) Color {
	return Color{clr: clr, val: val}
}

func (color Color) Print() string {
	return colorMap[color.clr](color.val)
}
