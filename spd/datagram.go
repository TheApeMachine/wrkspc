package spd

import (
	"bytes"
	"math/big"

	"github.com/theapemachine/wrkspc/errnie"
)

var (
	Version = []byte("v4")
	Empty   = Datagram{}
)

/*
New contructs a Datagram message and marshals it to a byte slice before returning
the object wrapped in a bytes.Buffer.
*/
func New(media MediaType, role RoleType, scope ScopeType) *Datagram {
	errnie.Trace()

	var (
		dg  *Datagram
		err error
	)

	// Get a root instance from Cap 'n Proto.
	if dg, err = root(); errnie.Handles(err) != nil {
		return dg
	}

	// Set the context header values to determine the way this
	// datagram should be processed.
	if errnie.Handles(dg.setContext(media, role, scope)) != nil {
		return dg
	}

	return dg
}

/*
Prefix generates the canonical key under which the datagram
can be found in the data lake.
*/
func (dg *Datagram) Prefix() *bytes.Buffer {
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
