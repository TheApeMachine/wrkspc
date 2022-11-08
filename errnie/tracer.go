package errnie

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/theapemachine/wrkspc/tui"
)

/*
Traces the current file, line number, and function of the line
this function is called from.
*/
func Traces() {
	if !ambctx.Tracing {
		// We're not tracing, bail.
		return
	}

	// Handover from this conveniance function to the actual
	// implementation method.
	ambctx.tracer.Inspect(false)
}

func Times(t time.Time) {
	if !ambctx.Tracing {
		// We're not tracing, bail.
		return
	}

	timeTrack(t)
}

/*
Tracer looks at the runtime environment and is able to extract and
display various information of the low-level program state.
*/
type Tracer struct{}

/*
NewTracer constructs a Tracer type and returns an instance of it.
*/
func NewTracer() Tracer {
	return Tracer{}
}

/*
Inspect the current runtime environment and display the file, line number,
and function call that the call to Traces was made from.
*/
func (tracer Tracer) Inspect(snippet bool) {
	// Collect stack data.
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)

	// Using a series of string manipulations we extract the filename,
	// line number, and current function call.
	frame, _ := runtime.CallersFrames(pc[:n]).Next()
	fchunks := strings.Split(frame.File, "/")
	fstr := strings.Join(fchunks[len(fchunks)-2:], "/")
	fnchunks := strings.Split(frame.Function, "/")
	fnstr := fnchunks[len(fnchunks)-1]

	// Add additional styling around the terminal output, then write
	// the final log line to the console.
	fmt.Println(
		tui.NewLabel(" TRACE ").Print(),
		tui.NewColor(
			"MUTE", time.Now().Format("2006-01-02 15:04:05.000000"),
		).Print(),
		tui.NewIcon("flag"),
		tui.NewColor("DARK", fstr).Print(),
		tui.NewColor("MUTE", "("+fnstr+")").Print(),
		tui.NewColor("HIGH", strconv.Itoa(frame.Line)).Print(),
	)

	if snippet {
		wd, _ := os.Getwd()
		buf, _ := os.OpenFile(wd+"/"+fstr, os.O_RDONLY, os.ModeAppend.Perm())
		scanner := bufio.NewScanner(buf)
		lastLine := 0

		var builder strings.Builder

		for scanner.Scan() {
			if lastLine >= frame.Line-5 && lastLine <= frame.Line+5 {
				suffix := ""
				if lastLine == frame.Line {
					suffix = tui.NewIcon("explode")
				}
				builder.WriteString(
					fmt.Sprintf("%d %s\n", lastLine, scanner.Text()+suffix),
				)
			}
			lastLine++
		}

		fmt.Println(
			tui.NewColor("DARK", builder.String()).Print(),
		)
	}
}

func timeTrack(start time.Time) {
	fmt.Println(tui.NewLabel("RUNTIME").Print(), time.Since(start).String())
}
