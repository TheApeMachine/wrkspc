package spd

import (
	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

func (datagram *Datagram) Decode(p []byte) error {
	var (
		msg *capnp.Message
		dg  Datagram
		err error
	)

	if msg, err = capnp.Unmarshal(p); errnie.Handles(err) != nil {
		return err
	}

	if dg, err = ReadRootDatagram(msg); errnie.Handles(err) != nil {
		return err
	}

	datagram = &dg
	return err
}

func (datagram *Datagram) Encode(p []byte) error {
	var (
		b   []byte
		err error
	)

	if b, err = datagram.Message().Marshal(); err != nil {
		return errnie.Handles(err)
	}

	copy(p, b)
	return err
}
