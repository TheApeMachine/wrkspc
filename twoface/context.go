package twoface

import (
	"context"
	"time"
)

type stringKey string

/*
Context is a conveniance wrapper around Go contexts to improve
developer ergonomics.
*/
type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
}

/*
NewContext constructs a twoface context.
*/
func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	deadline = time.Now().Add(3 * time.Second)
	ok = true
	return
}

func (ctx *Context) Done() <-chan struct{} {
	return make(chan struct{})
}

func (ctx *Context) Err() error {
	return nil
}

func (ctx *Context) Write(key stringKey, value any) {
	ctx.ctx = context.WithValue(ctx.ctx, key, value)
}

func (ctx *Context) Read(key stringKey) any {
	return ctx.ctx.Value(key)
}
