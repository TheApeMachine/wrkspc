package errnie

import (
	"runtime"
	"strings"

	"github.com/theapemachine/wrkspc/berrt"
)

type Tracer struct {
	diagram *berrt.Diagram
}

func Trace() {
	if !ctx.tracing {
		return
	}

	// Collect stack data.
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)

	// Using a series of string manipulations we extract the filename,
	// line number, and current function call.
	frame, _ := runtime.CallersFrames(pc[:n]).Next()
	fchunks := strings.Split(frame.File, "/")
	fstr := strings.Join(fchunks[len(fchunks)-2:], "/")
	fnchunks := strings.Split(frame.Function, "/")
	fnstr := fnchunks[len(fnchunks)-1]

	write("TRACE", fstr, fnstr, frame.Line)
}
