package spd

import (
	"bytes"
	"math/big"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
)

var Version = []byte("v4")

/*
New contructs a Datagram message and marshals it to a byte slice before returning
the object wrapped in a bytes.Buffer.
*/
func New(role, scope, identity []byte, layers []*bytes.Buffer) Datagram {
	errnie.Trace()

	arena := capnp.SingleSegment(nil)
	_, seg, err := capnp.NewMessage(arena)
	errnie.Handles(err)

	dg, err := NewRootDatagram(seg)
	errnie.Handles(err)

	errnie.Handles(
		dg.SetUuid([]byte(uuid.New().String())),
	)

	errnie.Handles(dg.SetVersion(Version))
	dg.SetTimestamp(time.Now().UnixNano())

	// Add a new layer to store our payload.
	list, err := capnp.NewDataList(seg, 1)
	errnie.Handles(err)

	// Write the layers to the Datagram.
	for idx, layer := range layers {
		errnie.Handles(list.Set(idx, layer.Bytes()))
	}
	errnie.Handles(dg.SetLayers(list))

	// Set the context header values to determine the way this
	// datagram should be processed.
	errnie.Handles(dg.SetRole(role))
	errnie.Handles(dg.SetScope(scope))
	errnie.Handles(dg.SetIdentity(identity))

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
