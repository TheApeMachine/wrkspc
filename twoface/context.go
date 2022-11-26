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
type Context interface {
	error
	io.ReadWriteCloser
	Get() context.Context
	Cancel() errnie.Error
	TTL() <-chan struct{}
}

/*
NewContext constructs a twoface.Context.
*/
func NewContext(contextType Context) Context {
	if contextType == nil {
		return NewProtoContext(10)
	}

	return contextType
}

/*
ProtoContext is the default context which should be able to handle
most use-cases in which you either need a Go context or a
twoface.Context wrapper.
*/
type ProtoContext struct {
	ctx context.Context
	cnl context.CancelFunc
	ttl context.CancelFunc
	err errnie.Error
	clk int
}

/*
NewProtoContext constructs a default context which should suffice
for most use-cases.
*/
func NewProtoContext(clk int) *ProtoContext {
	ctx := context.Background()
	ctx, cnl := context.WithCancel(ctx)
	ctx, ttl := context.WithTimeout(ctx, time.Duration(clk)*time.Second)

	return &ProtoContext{
		ctx: ctx,
		cnl: cnl,
		ttl: ttl,
		clk: clk,
	}
}

/*
Get the inner native Go context, so we can pass it when a function
or method expects to have one.
*/
func (proto *ProtoContext) Get() context.Context {
	return proto.ctx
}

/*
Cancel the current context and propegate the termination signal to
all goroutines that share the context.
*/
func (proto *ProtoContext) Cancel() errnie.Error {
	proto.cnl()
	return proto.err
}

/*
TTL returns a CancelFunc when the specified timeout has been reached.
*/
func (proto *ProtoContext) TTL() <-chan struct{} {
	return proto.ctx.Done()
}

/*
Error implements the Go error interface.
*/
func (proto *ProtoContext) Error() string {
	return proto.err.Error()
}

/*
Read implements the io.Reader interface.
*/
func (proto *ProtoContext) Read(p []byte) (n int, err error) {
	return
}

type stringKey string

/*
Write implements the io.Writer interface.
*/
func (proto *ProtoContext) Write(p []byte) (n int, err error) {
	return
}

/*
Close implements the io.Closer interface.
*/
func (proto *ProtoContext) Close() error {
	proto.Cancel()
	return errnie.NewError(nil)
}
