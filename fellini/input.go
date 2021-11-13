package fellini

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
NewInput constructs an input of the type that is passed in and calls its initialize method.
*/
func NewInput(
	inputType Component,
	channel chan []byte,
	disposer *twoface.Disposer,
) Component {
	return inputType.Initialize(channel, disposer)
}

/*
TextInput represents a text field on the terminal.
*/
type TextInput struct {
	channel     chan []byte
	disposer    *twoface.Disposer
	buf         string
	subtle      lipgloss.AdaptiveColor
	dialogStyle lipgloss.Style
}

/*
Initialize sets the Input up to be ready for use.
*/
func (component TextInput) Initialize(channel chan []byte, disposer *twoface.Disposer) Component {
	component.channel = channel
	component.disposer = disposer
	component.subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	component.dialogStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	// Listen to the channel for bytes coming in and buffer them in our Component type.
	go func() {
		for char := range component.channel {
			component.buf += string(char)

			prompt := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(fmt.Sprintf("amsh[>%s]", component.buf))
			dialog := lipgloss.Place(96, 9,
				lipgloss.Center, lipgloss.Center,
				component.dialogStyle.Render(prompt),
				lipgloss.WithWhitespaceChars("猫咪"),
				lipgloss.WithWhitespaceForeground(component.subtle),
			)

			fmt.Print(dialog)
		}
	}()

	return component
}
