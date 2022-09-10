package twoface

import (
	"context"
	"time"
)

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

func (ctx *Context) Value(key any) any {
	return nil
}
