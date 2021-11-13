package spdg

import (
	"bytes"
	"time"
)

/*
Datagram is a type that can wrap any kind of data and turn it into a common object that
can transport that data to any kind of method, function, channel, etc. that uses the
Datagram type. The idea is to use this in as many places as possible, so everything
speaks the same type. This in turn would result in a `mono typed` system, which has shown
some benefits and flexibilities that are otherwise not available.
*/
type Datagram struct {
	Context *Context
	Data    *Data
}

/*
NewDatagram constructs a Datagram shaped message by passing in preconfigured Context and Data
objects and returns a pointer reference to it.
*/
func NewDatagram(ctx *Context, dat *Data) *Datagram {
	return &Datagram{Context: ctx, Data: dat}
}

/*
QuickDatagram returns a Datagram that needs minimal additional configuration.
*/
func QuickDatagram(role ContextRole, dataType string, dat *bytes.Buffer) *Datagram {
	return NewDatagram(
		NewContext(role, time.Now().UnixNano(), dataType),
		NewData(NewHeader(), dat),
	)
}

/*
ContextDatagram returns a Datagram that only has a Context.
*/
func ContextDatagram(role ContextRole, annotations ...Annotation) *Datagram {
	return NewDatagram(
		&Context{
			Role:        role,
			Timestamp:   time.Now().UnixNano(),
			Annotations: annotations,
		}, nil,
	)
}
