package spd

import (
	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Decode a byte slice buffer into a Datagram type, so we can access its
fields and values.

This is a very cheap operation, because of how Cap 'n Proto works.
*/
func (datagram *Datagram) Decode(p []byte) error {
	var (
		msg *capnp.Message
		dg  Datagram
		err error
	)

	// Unmarshal is a bit of a misnomer in the world of Cap 'n Proto,
	// but they went with it anyway.
	if msg, err = capnp.Unmarshal(p); errnie.Handles(err) != nil {
		return err
	}

	// Read a Datagram instance from the message.
	if dg, err = ReadRootDatagram(msg); errnie.Handles(err) != nil {
		return err
	}

	// Overwrite the pointer to our empty instance with the one
	// pointing to our root Datagram.
	datagram = &dg
	return err
}

/*
Encode the Datagram to a byte slice buffer, so we are compatible with
lower-level operations (such as data storage).

This is a very cheap operation, because of how Cap 'n Proto works.
*/
func (datagram *Datagram) Encode() []byte {
	var (
		buf []byte
		err error
	)

	if buf, err = datagram.Message().Marshal(); err != nil {
		// Log the error, we don't need to return here, because we
		// rely on the caller handling things if the buffer does not
		// contain the expected data.
		errnie.Handles(err)
	}

	// Return the buffer in whatever state it may be.
	return buf
}
