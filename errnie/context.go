package errnie

import (
	"io"
	"os"

	"github.com/theapemachine/wrkspc/berrt"
)

/*
ctx captures the internal state and behavior of errnie as an
ambient context, which can be accessed anywhere through the
publicly exposed functions.
*/
var ctx *Context

/*
init makes sure that the ambient context is loaded up and instantiated
before any other application code executes.
*/
func init() {
	ctx = New()
}

/*
Context wraps all state and behavior errnie needs to act as an
error handler, logger, and tracer.
*/
type Context struct {
	output      io.Writer
	fh          *os.File
	tracing     bool
	debugging   bool
	breakpoints bool
}

/*
New constructs and instantiates the ambient context available to errnie
internally. Application code can access the instance through the
publicly exposed functions.
*/
func New() *Context {
	wd, err := os.Getwd()
	Handles(err)

	fh, err := os.Open(wd + "/log")
	Handles(err)

	// Return the context instance loaded with any desired output
	// channels (of type io.Writer). Our logging operations will
	// pipe data directly to them.
	return &Context{
		output: io.MultiWriter(
			os.Stdout,
			fh,
			berrt.NewSequenceDiagram(),
		),
		fh: fh,
	}
}

func Quiet(output io.Writer) {
	ctx.output = io.MultiWriter(output)
}

/*
Ctx returns the ambient context for use of its io.ReadWriteCloser
interface implementation.
*/
func Ctx() *Context { return ctx }

/* Tracing behavior turned on or off. */
func Tracing(value bool) { ctx.tracing = value }

/* Debugging behavior turned on or off. */
func Debugging(value bool) { ctx.debugging = value }

/* Breakpoints behavior turned on or off. */
func Breakpoints(value bool) { ctx.breakpoints = value }
