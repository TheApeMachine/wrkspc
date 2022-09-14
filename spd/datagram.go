package spd

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
)

type cacheTuple struct {
	dg  *Datagram
	seg *capnp.Segment
	msg *capnp.Message
}

var pool = sync.Pool{
	New: func() interface{} {
		arena := capnp.SingleSegment(nil)
		msg, seg, err := capnp.NewMessage(arena)
		errnie.Handles(err).With(errnie.NOOP)

		dg, err := NewRootDatagram(seg)
		errnie.Handles(err).With(errnie.NOOP)

		errnie.Handles(dg.SetUuid(uuid.NewString())).With(errnie.NOOP)
		errnie.Handles(dg.SetVersion("v4.0.0")).With(errnie.NOOP)
		dg.SetTimestamp(time.Now().UnixNano())

		return &cacheTuple{dg: &dg, seg: seg, msg: msg}
	},
}

func NewCached(role, scope, identity, payload string) []byte {
	errnie.Traces()

	// Retrieve a pre-made datagram from the cache, which
	// already have a version, timestamp, and uuid.
	cTup := pool.Get().(*cacheTuple)
	defer pool.Put(cTup)

	// Add a new layer to store our payload.
	list, err := capnp.NewDataList(cTup.seg, 1)
	errnie.Handles(err).With(errnie.NOOP)
	errnie.Handles(list.Set(0, []byte(payload))).With(errnie.NOOP)
	errnie.Handles(cTup.dg.SetLayers(list)).With(errnie.NOOP)

	// Set the context header values to determine the way this
	// datagram should be processed.
	errnie.Handles(cTup.dg.SetRole(role)).With(errnie.NOOP)
	errnie.Handles(cTup.dg.SetScope(scope)).With(errnie.NOOP)
	errnie.Handles(cTup.dg.SetIdentity(identity)).With(errnie.NOOP)

	// Marshal into a byte slice so it becomes compatible with
	// network and storage layers.
	b, err := cTup.msg.Marshal()
	errnie.Handles(err).With(errnie.NOOP)

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
	errnie.Handles(err).With(errnie.NOOP)

	role, err := dg.Role()
	errnie.Handles(err).With(errnie.NOOP)

	scope, err := dg.Scope()
	errnie.Handles(err).With(errnie.NOOP)

	identity, err := dg.Identity()
	errnie.Handles(err).With(errnie.NOOP)

	timestamp := dg.Timestamp()

	uuid, err := dg.Uuid()
	errnie.Handles(err).With(errnie.NOOP)

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
	errnie.Handles(err).With(errnie.NOOP)

	data, err := list.At(0)
	errnie.Handles(err).With(errnie.NOOP)

	return data
}

func Unmarshal(p []byte) Datagram {
	msg, err := capnp.Unmarshal(p)
	errnie.Handles(err).With(errnie.NOOP)

	dg, err := ReadRootDatagram(msg)
	errnie.Handles(err).With(errnie.NOOP)

	return dg
}
