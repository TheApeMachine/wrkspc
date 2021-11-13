package errnie

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/berrt"
)

/*
Tracer is an object that inspects the stack.
*/
type Tracer struct {
	labelStyle     func(string) string
	mutedStyle     func(string) string
	darkStyle      func(string) string
	normalStyle    func(string) string
	highlightStyle func(string) string
}

/*
NewTracer constructs a Tracer and returns a reference pointere to it.
*/
func NewTracer() *Tracer {
	return &Tracer{
		labelStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#626262"),
		).Background(
			lipgloss.Color("#262626"),
		).Bold(true).Padding(0, 1).Render,

		mutedStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#A8A8A8"),
		).Render,

		darkStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#808080"),
		).Align(lipgloss.Left).Render,

		normalStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#444444"),
		).Render,

		highlightStyle: lipgloss.NewStyle().Foreground(
			lipgloss.Color("#EEEEEE"),
		).Render,
	}
}

/*
Runtime retrieves runtime information periodically and outputs this to the terminal.
It takes an interval (in seconds) to determine the time between updates.
*/
func (tracer *Tracer) Runtime(interval int) {
	if !viper.GetBool("wrkspc.errnie.trace") {
		return
	}

	for {
		fmt.Printf(
			"%s %s\n",
			berrt.NewLabel("RUNTIME").ToString(),
			berrt.NewText(fmt.Sprintf("GOROUTINES: %v", runtime.NumGoroutine())).ToString(),
		)

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

/*
Inspect the Stack and print the tracing output.
*/
func (tracer *Tracer) Inspect(flags ...bool) {
	if !viper.GetBool("wrkspc.errnie.trace") {
		return
	}

	// Collect stack data.
	pc := make([]uintptr, 15)
	var n int

	if len(flags) == 2 && flags[1] {
		n = runtime.Callers(7, pc)
	} else {
		n = runtime.Callers(3, pc)
	}

	frame, _ := runtime.CallersFrames(pc[:n]).Next()

	// Format and shape into strings.
	fchunks := strings.Split(frame.File, "/")
	fstr := strings.Join(fchunks[len(fchunks)-2:], "/")
	fnchunks := strings.Split(frame.Function, "/")
	fnstr := fnchunks[len(fnchunks)-1]
	tstr := time.Now().Format("2006-01-02 15:04:05.000000")

	// Style the strings.
	label := tracer.labelStyle(" TRACE ")
	fStyled := tracer.mutedStyle(fstr)
	fnStyled := tracer.darkStyle("(" + fnstr + ")")
	tStyled := tracer.normalStyle(tstr)
	lStyled := tracer.highlightStyle(strconv.Itoa(frame.Line))

	icon := "\xF0\x9F\x94\xB9"

	if strings.Split(fstr, "/")[0] == "errnie" {
		icon = "\xF0\x9F\x94\xB8"
	}

	fmt.Printf("%s %s %s %s %s %s\n", label, tStyled, icon, fStyled, lStyled, fnStyled)
	tracer.renderCode(frame.File, frame.Line, flags...)
}

/*
renderCode takes the file name and line number and prints the code that is found
using those values. It performs a check to see if the `code` param is nil or false,
such that we can short circuit.
*/
func (tracer *Tracer) renderCode(fp string, ln int, flags ...bool) {
	if len(flags) == 0 || !flags[0] {
		return
	}

	var out string

	r, _ := os.Open(fp)
	lastline := 0
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lastline++

		if lastline >= ln-5 && lastline <= ln+5 {
			out += tracer.setStyle(scanner, lastline, ln) + "\n"
		}
	}

	fmt.Printf("%s", out)
}

/*
setStyle was extracted because I hate nesting levels that go too deep, but really its an
indication that a new type should exist for this as it looks out of place and it is.
*/
func (tracer *Tracer) setStyle(scanner *bufio.Scanner, lastline int, ln int) string {
	style := tracer.darkStyle
	icon := ""

	if lastline == ln {
		style = tracer.highlightStyle
		icon = "\xF0\x9F\x94\xB4"
	}

	return style(strconv.Itoa(lastline) + " " + scanner.Text() + " " + icon)
}
