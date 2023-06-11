package spd

import (
	"bytes"

	"capnproto.org/go/capnp/v3"
	"github.com/wrk-grp/errnie"
)

var Version = []byte("0.0.5")
var Empty = &Datagram{}

func root() (*Datagram, error) {
	errnie.Trace()

	arena := capnp.SingleSegment(nil)
	_, seg, err := capnp.NewMessage(arena)

	if errnie.Handles(err) != nil {
		return Empty, err
	}

	var dg Datagram

	if dg, err = NewRootDatagram(seg); err != nil {
		return Empty, errnie.Handles(err)
	}

	return &dg, nil
}

/*
Prefix uses the context header values to generate a prefix for the message.
The canonical prefix is made out of the following values:

  - Version
  - Type
  - Role
  - Scope
  - Identity
  - Uuid

The prefix is used to determine the path to the file in S3.
*/
func (dg *Datagram) Prefix() []byte {
	errnie.Trace()

	var (
		err error
	)

	t, err := dg.Type()
	errnie.Handles(err)

	r, err := dg.Role()
	errnie.Handles(err)

	s, err := dg.Scope()
	errnie.Handles(err)

	i, err := dg.Identity()
	errnie.Handles(err)

	u, err := dg.Uuid()
	errnie.Handles(err)

	var buf bytes.Buffer
	buf.Write(Version)
	buf.WriteByte('/')
	buf.Write(t)
	buf.WriteByte('/')
	buf.Write(r)
	buf.WriteByte('/')
	buf.Write(s)
	buf.WriteByte('/')
	buf.Write(i)
	buf.WriteByte('/')
	buf.Write(u)

	return buf.Bytes()
}
