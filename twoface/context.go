package twoface

import (
	"context"
	"io"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Context is a wrapper around native the Go context type,
while adding functionality to improve developer ergonomics.
*/
type Context struct {
	error
	io.ReadWriteCloser
	context.Context

	root   context.Context
	cancel context.CancelFunc
	ttl    context.CancelFunc

	data []byte
	rIdx int64
	err  error
}

/*
NewContext constructs a twoface.Context.
*/
func NewContext() *Context {
	// Construct the underlying native Go Context with a CancelFunc
	// and Deadline/Timeout.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx, ttl := context.WithTimeout(ctx, 10*time.Second)

	return &Context{
		root:   ctx,
		cancel: cancel,
		ttl:    ttl,
		data:   make([]byte, 0),
	}
}

/*
Error implements Go's native error interface and wraps errnie around
it to provide more context and more flexible output.
*/
func (ctx *Context) Error() string {
	return errnie.NewError(ctx.Err()).Error()
}

/*
Read implements the io.Reader interface.
*/
func (ctx *Context) Read(p []byte) (n int, err error) {
	if ctx.rIdx >= int64(len(ctx.data)) {
		err = io.EOF
		return
	}

	n = copy(p, ctx.data[ctx.rIdx:])
	ctx.rIdx += int64(n)
	return
}

/*
Write implements the io.Writer interface.
*/
func (ctx *Context) Write(p []byte) (n int, err error) {
	ctx.data = append(ctx.data, p...)
	return len(p), nil
}

/*
Close implements the io.Closer interface.
*/
func (ctx *Context) Close() error {
	ctx.cancel()
	return ctx.Err()
}

/*
Deadline ...
*/
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.root.Deadline()
}

/*
Done ...
*/
func (ctx *Context) Done() <-chan struct{} {
	return make(<-chan struct{})
}

/*
Err ...
*/
func (ctx *Context) Err() error {
	return errnie.NewError(ctx.err)
}

/*
Value ...
*/
func (ctx *Context) Value(key any) any {
	return key
}
