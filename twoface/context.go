package twoface

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

/*
Context is a wrapper around native the Go context type,
while adding functionality to improve developer ergonomics.
*/
type Context interface {
	error
	io.ReadWriter
	Get() context.Context
	Cancel() errnie.Error
	TTL() <-chan struct{}
}

/*
NewContext constructs a twoface.Context.
*/
func NewContext(contextType Context) Context {
	if contextType == nil {
		return NewProtoContext()
	}

	return contextType
}

/*
ProtoContext is the default context which should be able
to handle most use-cases in which you either need a Go
context or a twoface.Context wrapper.
*/
type ProtoContext struct {
	ctx context.Context
	cnl context.CancelFunc
	ttl context.CancelFunc
	err errnie.Error
}

/*
NewProtoContext constructs a default context which should
suffice for most use-cases.
*/
func NewProtoContext() *ProtoContext {
	ctx := context.Background()
	ctx, cnl := context.WithCancel(ctx)
	ctx, ttl := context.WithTimeout(ctx, time.Second*10)

	return &ProtoContext{
		ctx: ctx,
		cnl: cnl,
		ttl: ttl,
	}
}

/*
Get the inner native Go context, so we can pass it when
a function or method expects to have one.
*/
func (proto *ProtoContext) Get() context.Context {
	return proto.ctx
}

/*
Cancel the current context and propegate the termination
signal to all goroutines that share the context.
*/
func (proto *ProtoContext) Cancel() errnie.Error {
	proto.cnl()
	return proto.err
}

/*
TTL returns a CancelFunc when the specified timeout
has been reached.
*/
func (proto *ProtoContext) TTL() <-chan struct{} {
	return proto.ctx.Done()
}

/*
Error implements the Go error interface.
*/
func (proto *ProtoContext) Error() string {
	return proto.err.Msg
}

/*
Read implements the io.Reader interface.
*/
func (proto *ProtoContext) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, errors.New("no key")
	}

	var val interface{}
	key := string(p)

	if val = proto.ctx.Value(key); val == nil {
		return 0, errors.New("no value")
	}

	var ok bool

	if p, ok = val.([]byte); !ok {
		return 0, errors.New("type error")
	}

	return len(p), nil
}

type stringKey string

/*
Write implements the io.Writer interface.
*/
func (proto *ProtoContext) Write(p []byte) (n int, err error) {
	dg := spd.Unmarshal(p)
	proto.ctx = context.WithValue(
		proto.ctx, stringKey(dg.Prefix()), p,
	)
	return len(p), nil
}
