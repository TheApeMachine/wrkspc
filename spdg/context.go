package spdg

import (
	"github.com/mitchellh/hashstructure/v2"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
ContextRole wraps the role field into an enum.
*/
type ContextRole string

const (
	// NULLGRAM is used when the return value is void basically.
	NULLGRAM ContextRole = "null"
	// ANONYMOUS is used as a `generic` type.
	ANONYMOUS ContextRole = "anonymous"
	// ERROR is used as a `generic` error type.
	ERROR ContextRole = "error"
	// DATAPOINT is usually something you want to store somehow.
	DATAPOINT ContextRole = "datapoint"
	// COMMAND is to trigger processes in the data pipelines.
	COMMAND ContextRole = "command"
	// QUESTION is to search for data.
	QUESTION ContextRole = "question"
	// TOPIC is a key that targets data at a Question.
	TOPIC ContextRole = "topic"
)

/*
Context is a metadata header wrapping the Datagram. It describes
the payload such that the inner data can be abstracted away as
anonymous bytes.
*/
type Context struct {
	ID          uint64
	Role        ContextRole
	Timestamp   int64
	Type        string
	Annotations []Annotation
}

/*
NewContext constructs a Header for a Datagram, which contains the meta data
for the Payload inside the Body of the Datagram.
*/
func NewContext(role ContextRole, timestamp int64, dataType string) *Context {
	ctx := &Context{
		Role:      role,
		Timestamp: timestamp,
		Type:      dataType,
	}

	hash, err := hashstructure.Hash(ctx, hashstructure.FormatV2, nil)
	ctx.ID = hash

	errnie.Handles(err).With(errnie.NOOP)
	return ctx
}

/*
Prefix collapses the Context header into a path style prefix to be compatible with most cloud
storage engines, as well as radix trees.
*/
func (context *Context) Prefix() *string {
	str := ""
	return &str
}

/*
Annotate the Context which act as a flexible `parameter` to a use-case. Basically use these scoped
to your needs, in any way shape or form that is compatible.
*/
func (ctx *Context) Annotate(key, value string) *Context {
	ctx.Annotations = append(ctx.Annotations, NewAnnotation(key, value))
	return ctx
}
