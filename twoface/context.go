package twoface

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Context is a wrapper around native the Go context type,
while adding functionality to improve developer ergonomics.
*/
type Context struct {
	root   context.Context
	cancel context.CancelFunc
	ttl    context.CancelFunc
	wgs    []*sync.WaitGroup
	data   []byte
	rIdx   int64
	err    error
}

/*
NewContext constructs a twoface.Context.
*/
func NewContext() *Context {
	errnie.Trace()

	// Construct the underlying native Go Context with a CancelFunc
	// and Deadline/Timeout.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx, ttl := context.WithTimeout(ctx, 10*time.Second)

	return &Context{
		root:   ctx,
		cancel: cancel,
		ttl:    ttl,
		wgs:    make([]*sync.WaitGroup, 0),
		data:   make([]byte, 0),
	}
}

func (ctx *Context) Root() context.Context {
	errnie.Trace()
	return ctx.root
}

func (ctx *Context) WG(idx int, val int) *sync.WaitGroup {
	errnie.Trace()

	if idx > len(ctx.wgs)-1 {
		ctx.wgs = append(ctx.wgs, &sync.WaitGroup{})
	}

	if val == 1 {
		errnie.Debugs("wg", idx, val)
		ctx.wgs[idx].Add(1)
	}

	if val == -1 {
		errnie.Debugs("wg", idx, val)
		ctx.wgs[idx].Done()
	}

	return ctx.wgs[idx]
}

func (ctx *Context) Wait(idx int) (err error) {
	errnie.Trace()

	if idx > len(ctx.wgs)-1 {
		return errnie.Handles(errors.Wrap(err, "invalid waitgroup"))
	}

	ctx.wgs[idx].Wait()
	errnie.Success("done")

	return
}

/*
Error implements Go's native error interface and wraps errnie around
it to provide more context and more flexible output.
*/
func (ctx *Context) Error() string {
	errnie.Trace()
	return errnie.NewError(ctx.Err()).Error()
}

/*
Read implements the io.Reader interface.
*/
func (ctx *Context) Read(p []byte) (n int, err error) {
	errnie.Trace()

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
	errnie.Trace()
	ctx.data = append(ctx.data, p...)
	return len(p), nil
}

/*
Close implements the io.Closer interface.
*/
func (ctx *Context) Close() error {
	errnie.Trace()
	ctx.cancel()
	return ctx.Err()
}

/*
Deadline ...
*/
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	errnie.Trace()
	return ctx.root.Deadline()
}

/*
Done ...
*/
func (ctx *Context) Done() <-chan struct{} {
	errnie.Trace()
	return make(<-chan struct{})
}

/*
Err ...
*/
func (ctx *Context) Err() error {
	errnie.Trace()
	return errnie.NewError(ctx.err)
}

/*
Value ...
*/
func (ctx *Context) Value(key any, val any) any {
	errnie.Trace()
	ctx.root = context.WithValue(ctx.root, key, val)
	return key
}
