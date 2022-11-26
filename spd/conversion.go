package spd

import (
	"bytes"
	"io"

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
		buf := bytes.NewBuffer([]byte{})
		io.Copy(buf, err)
		dg.Write(buf.Bytes())
	}

	dg, err = ReadRootDatagram(msg)
	errnie.Handles(err)

	return dg
}
