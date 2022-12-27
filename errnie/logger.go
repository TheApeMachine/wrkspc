package errnie

import (
	"bytes"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type LogLevel uint

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
)

/*
Logger is an interface objects can implement if they want to act as
log output channels for errnie.
*/
type Logger interface {
	Error(...any)
	Warning(...any)
	Info(...any)
	Debug(...any)
	Inspect(...any)
}

var styles = map[string]lipgloss.Style{
	"darkest": lipgloss.NewStyle().Bold(true).Italic(true).Foreground(
		lipgloss.Color("#323232"),
	),
	"darker": lipgloss.NewStyle().Bold(true).Italic(true).Foreground(
		lipgloss.Color("#424242"),
	),
	"dark": lipgloss.NewStyle().Bold(true).Italic(true).Foreground(
		lipgloss.Color("#626262"),
	),
	"TRACE": lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.Color("#A2A2A2"),
	).Background(
		lipgloss.Color("#626262"),
	),
	"DEBUG": lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.Color("#626262"),
	).Background(
		lipgloss.Color("#A2A2A2"),
	),
	"INFO": lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.Color("#EFEFEF"),
	).Background(
		lipgloss.Color("#33AAFF"),
	),
	"SUCCESS": lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.Color("#EFEFEF"),
	).Background(
		lipgloss.Color("#00FF55"),
	),
	"WARNING": lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.Color("#EFEFEF"),
	).Background(
		lipgloss.Color("#FFAA33"),
	),
	"ERROR": lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.Color("#EFEFEF"),
	).Background(
		lipgloss.Color("#FF0055"),
	),
}

func write(level string, msgs ...any) {
	buf := bytes.NewBuffer([]byte{})

	t := time.Now().Format("2006-01-02 15:04:05.999999")
	for len(t) < 26 {
		t += "0"
	}

	buf.WriteString(
		styles["darkest"].Render("[") +
			styles["dark"].Render(t) +
			styles["darkest"].Render("]"),
	)
	buf.WriteString(" ")

	prfx, sufx := "", ""

	for len(prfx+level+sufx) < 7 {
		if len(prfx+level+sufx)%2 == 1 {
			prfx += " "
			continue
		}

		sufx += " "
	}

	buf.WriteString(styles[level].Render(prfx + level + sufx))

	for _, msg := range msgs {
		buf.WriteString(" ")
		fmt.Fprintf(buf, "%v", msg)
	}

	buf.WriteString("\n")
	ctx.Write(buf.Bytes())
}

/*
Errors is syntactic sugar to call the Warning method on
a Logger interface.
*/
func Errors(msgs ...any) {
	write("ERROR", msgs...)
}

/*
Warns is syntactic sugar to call the Warning method on
a Logger interface.
*/
func Warns(msgs ...any) {
	write("WARNING", msgs...)
}

/*
Success is syntactic sugar to call the Info method on
a Logger interface.
*/
func Success(msgs ...any) {
	write("SUCCESS", msgs...)
}

/*
Informs is syntactic sugar to call the Info method on
a Logger interface.
*/
func Informs(msgs ...any) {
	write("INFO", msgs...)
}

/*
Debugs is syntactic sugar to call the Debug method on
a Logger interface.
*/
func Debugs(msgs ...any) []byte {
	if ctx.debugging {
		write("DEBUG", msgs...)
	}
	return []byte(fmt.Sprint(msgs...))
}

/*
Inspects is syntactic sugar to dump the structure and values
of objects with arbitrary complexity to logger output channels.
*/
func Inspects(msgs ...any) {
	buf := bytes.NewBuffer([]byte{})

	for _, msg := range msgs {
		spew.Fprintf(buf, "%v", msg)
	}

	buf.WriteString("\n")
	ctx.Write(buf.Bytes())
}
