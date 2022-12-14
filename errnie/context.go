package errnie

import (
	"io"
	"os"

	"github.com/pkg/errors"
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
	go func() {
		defer ctx.reader.Close()

		for {
			if _, err := io.Copy(ctx.log, ctx.reader); err != nil {
				Handles(err)
				Handles(ctx.reader.CloseWithError(
					errors.Wrap(err, "pipe closed with error"),
				))
				return
			}
		}
	}()
}

/*
Context wraps all state and behavior errnie needs to act as an
error handler, logger, and tracer.
*/
type Context struct {
	log         io.Writer
	reader      *io.PipeReader
	writer      *io.PipeWriter
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

	f, err := os.Open(wd + "/log")
	Handles(err)

	r, w := io.Pipe()

	return &Context{
		log:    io.MultiWriter(os.Stdout, f),
		reader: r,
		writer: w,
	}
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
