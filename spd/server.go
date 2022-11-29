package spd

import (
	context "context"

	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

type AnnounceServer struct{}

func (AnnounceServer) Public(ctx context.Context, call Announce_public) error {
	var res Announce_public_Results
	var err error
	var lst capnp.DataList
	var msg *capnp.Message

	if res, err = call.AllocResults(); errnie.IOError(err) != nil {
		return err
	}

	// Set the result to be the product of the two arguments, A and B,
	// that we received. These are found in the Arith_multiply struct.
	var buf []byte
	if buf, err = call.Args().Datagram(); errnie.ConversionError(err) != nil {
		return err
	}

	var dg Datagram
	if msg, err = capnp.Unmarshal(buf); err != nil {
		if dg, err = ReadRootDatagram(msg); err != nil {
			// Do something with the Datagram.
			errnie.Inspects(dg)
		}
	}

	res.SetCrowd(lst)
	return nil
}
