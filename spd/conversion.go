package spd

import (
	capnp "capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Unmarshal the byte slice into a Datagram. This actually does
no desrialization at all, given Cap 'n Proto is operating
directly on byte arrays.
*/
func (dg Datagram) Unmarshal(p []byte) Datagram {
	msg, err := capnp.Unmarshal(p)

	if err := errnie.Handles(err); err != nil {
		l, ee := dg.Layers()
		if errnie.Handles(ee) != nil {
			return dg
		}

		var p []byte
		err.Write(p)
		l.Set(l.Len(), p)
	}

	dg, err = ReadRootDatagram(msg)
	errnie.Handles(err)

	return dg
}
