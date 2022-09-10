package spd

import (
	"strconv"
	"strings"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
)

var cache chan cacheTuple
var kill chan struct{}

type cacheTuple struct {
	dg  *Datagram
	msg *capnp.Message
}

func InitCache() {
	cache = make(chan cacheTuple, 128)

	go func() {
		defer close(cache)

		for {
			select {
			case <-kill:
				return
			default:
				arena := capnp.SingleSegment(nil)
				msg, seg, err := capnp.NewMessage(arena)
				errnie.Handles(err).With(errnie.NOOP)

				dg, err := NewRootDatagram(seg)
				errnie.Handles(err).With(errnie.NOOP)

				errnie.Handles(dg.SetUuid(uuid.NewString())).With(errnie.NOOP)
				errnie.Handles(dg.SetVersion("v4.0.0")).With(errnie.NOOP)
				dg.SetTimestamp(time.Now().UnixNano())

				cache <- cacheTuple{dg: &dg, msg: msg}
			}
		}
	}()
}

func NewCached(role, scope, identity string) []byte {
	errnie.Traces()

	cTup := <-cache

	errnie.Handles(cTup.dg.SetRole(role)).With(errnie.NOOP)
	errnie.Handles(cTup.dg.SetScope(scope)).With(errnie.NOOP)
	errnie.Handles(cTup.dg.SetIdentity(identity)).With(errnie.NOOP)

	b, err := cTup.msg.Marshal()
	errnie.Handles(err).With(errnie.NOOP)

	return b
}

func Prefix(dg Datagram) string {
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
