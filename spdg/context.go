package spdg

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Context is a metadata header wrapping the Datagram. It describes
the payload such that the inner data can be abstracted away as
anonymous bytes.
*/
type Context struct {
	prefix      *string
	ID          uuid.UUID    `json:"id"`
	Role        ContextRole  `json:"role"`
	Timestamp   int64        `json:"timestamp"`
	Type        string       `json:"type"`
	Annotations []Annotation `json:"annotations"`
}

/*
NewContext constructs a Header for a Datagram, which contains the meta data
for the Payload inside the Body of the Datagram.
*/
func NewContext(role ContextRole, timestamp int64, dataType string) *Context {
	id, err := uuid.NewUUID()
	errnie.Handles(err).With(errnie.NOOP)

	ctx := &Context{
		Role:      role,
		Timestamp: timestamp,
		Type:      dataType,
		ID:        id,
	}

	return ctx
}

/*
Prefix collapses the Context header into a path style prefix to be compatible with most cloud
storage engines, as well as radix trees.
*/
func (context *Context) Prefix() *string {
	// Only write it if we don't have a prefix already. Once a prefix
	// is set, it should be immutable!
	if context.prefix != nil {
		return context.prefix
	}

	var builder strings.Builder

	context.writeContextAnnotations(&builder)

	builder.WriteString(string(context.Role))
	builder.WriteString("/")
	builder.WriteString(strings.Join(
		strings.FieldsFunc(
			time.Unix(0, context.Timestamp).String(), Split,
		)[:6], "/",
	))
	builder.WriteString("/")
	builder.WriteString(context.ID.String())
	builder.WriteString(".json")

	out := builder.String()
	context.prefix = &out

	return context.prefix
}

/*
writeContextAnnotations searches for special annotations to add to the prefix.
*/
func (ctx *Context) writeContextAnnotations(builder *strings.Builder) {
	collector := make([]string, 2)

	for _, annotation := range ctx.Annotations {
		if annotation.Key == "canon" {
			collector[0] = annotation.Value
		}

		if annotation.Key == "identity" {
			collector[1] = annotation.Value
		}
	}

	for _, collected := range collector {
		builder.WriteString(collected)
		builder.WriteString("/")
	}
}

/*
Annotate the Context which act as a flexible `parameter` to a use-case. Basically use these scoped
to your needs, in any way shape or form that is compatible.
*/
func (ctx *Context) Annotate(key, value string) *Context {
	ctx.Annotations = append(ctx.Annotations, NewAnnotation(key, value))
	return ctx
}

/*
Validate the Context of the Datagram to verify state is correct before performing any operation on, or with
this instance. Unvalidated Datagrams are only supposed to travel over dumb pipes.
*/
func (ctx *Context) Validate() ContextError {
	if err := Check(ctx.Type, "string").IsNot("unk"); err != CONTEXTOK {
		return err
	}

	if err := Check(ctx.Role, "ContextRole").IsNot(BASEGRAM); err != CONTEXTOK {
		return err
	}

	if err := Check(ctx, "").IsComplete(); err != CONTEXTOK {
		return err
	}

	return CONTEXTOK
}

/*
Split is a custom comparison method for splitting strings so we can split a UTC timestamp format on
all the delimiters that are present in it. This method is still needed for Prefix generation
even now we use UnixNano timestamps, because we convert those to UTC to build the Prefix.
*/
func Split(delim rune) bool {
	return delim == '-' || delim == ':' || delim == ' ' || delim == '.'
}
