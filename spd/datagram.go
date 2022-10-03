package spd

import (
	"strconv"
	"strings"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
)

func NewCached(role, scope, identity, payload string) []byte {
	errnie.Traces()

	arena := capnp.SingleSegment(nil)
	msg, seg, err := capnp.NewMessage(arena)
	errnie.Handles(err)

	dg, err := NewRootDatagram(seg)
	errnie.Handles(err)

	errnie.Handles(dg.SetUuid(uuid.NewString()))
	errnie.Handles(dg.SetVersion("v4.0.0"))
	dg.SetTimestamp(time.Now().UnixNano())

	// Add a new layer to store our payload.
	list, err := capnp.NewDataList(seg, 1)
	errnie.Handles(err)
	errnie.Handles(list.Set(0, []byte(payload)))
	errnie.Handles(dg.SetLayers(list))

	// Set the context header values to determine the way this
	// datagram should be processed.
	errnie.Handles(dg.SetRole(role))
	errnie.Handles(dg.SetScope(scope))
	errnie.Handles(dg.SetIdentity(identity))

	// Marshal into a byte slice so it becomes compatible with
	// network and storage layers.
	b, err := msg.Marshal()
	errnie.Handles(err)

	return b
}

/*
Prefix generates the canonical key under which the datagram
can be found in the data lake.
*/
func (dg Datagram) Prefix() string {
	errnie.Traces()

	var builder strings.Builder

	version, err := dg.Version()
	errnie.Handles(err)

	role, err := dg.Role()
	errnie.Handles(err)

	scope, err := dg.Scope()
	errnie.Handles(err)

	identity, err := dg.Identity()
	errnie.Handles(err)

	timestamp := dg.Timestamp()

	uuid, err := dg.Uuid()
	errnie.Handles(err)

	builder.WriteString(version)
	builder.WriteString("/")
	builder.WriteString(role)
	builder.WriteString("/")
	builder.WriteString(scope)
	builder.WriteString("/")
	builder.WriteString(identity)
	builder.WriteString("/")
	builder.WriteString(strconv.FormatInt(timestamp, 10))
	builder.WriteString("/")
	builder.WriteString(uuid)

	return builder.String()
}

func Payload(dg Datagram) []byte {
	list, err := dg.Layers()
	errnie.Handles(err)

	data, err := list.At(0)
	errnie.Handles(err)

	return data
}

/*
Unmarshal the byte slice into a Datagram. This actually does
no desrialization at all, given Cap 'n Proto is operating
directly on byte arrays.
*/
func Unmarshal(p []byte) Datagram {
	msg, err := capnp.Unmarshal(p)

	if err := errnie.Handles(err); err.Type != errnie.NIL {
		m, e := capnp.Unmarshal(
			NewCached("error", "unmarshal", "wrkspc", err.Msg),
		)

		errnie.Handles(e)

		dg, e := ReadRootDatagram(m)
		errnie.Handles(e)
		return dg
	}

	dg, err := ReadRootDatagram(msg)
	errnie.Handles(err)

	return dg
}
