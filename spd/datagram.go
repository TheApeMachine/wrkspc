package spd

import (
	"bytes"
	"math/big"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
)

var (
	Version = []byte("v4")
	Empty   = Datagram{}
)

/*
New contructs a Datagram message and marshals it to a byte slice before returning
the object wrapped in a bytes.Buffer.
*/
func New(
	role RoleType,
	scope ScopeType,
) Datagram {
	errnie.Trace()

	arena := capnp.SingleSegment(nil)
	_, seg, err := capnp.NewMessage(arena)
	if errnie.Handles(err) != nil {
		return Empty
	}

	dg, err := NewRootDatagram(seg)
	if errnie.Handles(err) != nil {
		return Empty
	}

	if errnie.Handles(
		dg.SetUuid([]byte(uuid.New().String())),
	) != nil {
		return Empty
	}

	if errnie.Handles(dg.SetVersion(Version)) != nil {
		return Empty
	}

	dg.SetTimestamp(time.Now().UnixNano())

	// Set the context header values to determine the way this
	// datagram should be processed.
	if errnie.Handles(dg.SetRole(role)) != nil {
		return Empty
	}

	if errnie.Handles(dg.SetScope(scope)) != nil {
		return Empty
	}

	if errnie.Handles(dg.SetIdentity(tweaker.GetIdentity())) != nil {
		return Empty
	}

	return dg
}

/*
Prefix generates the canonical key under which the datagram
can be found in the data lake.
*/
func (dg Datagram) Prefix() *bytes.Buffer {
	errnie.Trace()

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

	builder := bytes.NewBuffer([]byte{})
	for _, i := range [][]byte{version, role, scope, identity} {
		builder.Write(i)
		builder.WriteString("/")
	}

	big := new(big.Int)
	big.SetInt64(timestamp)

	builder.Write(big.Bytes())
	builder.WriteString("/")
	builder.Write(uuid)

	return builder
}
