package spd

import (
	"bytes"
	context "context"
	"log"
	"net"

	"capnproto.org/go/capnp/v3/rpc"
	"github.com/theapemachine/wrkspc/errnie"
)

type AnnounceClient struct{}

func NewAnnounceClient() *AnnounceClient {
	network := "tcp"
	connection, err := net.Dial(network, "127.0.0.1")
	errnie.Handles(err)

	ctx := context.Background()
	conn := rpc.NewConn(rpc.NewStreamTransport(connection), nil)
	defer conn.Close()

	announce := Announce(conn.Bootstrap(ctx))
	datagram := New(
		[]byte("test"), []byte("test"), []byte("test"),
		[]*bytes.Buffer{},
	)

	f, release := announce.Public(ctx, func(ps Announce_public_Params) error {
		dg, err := datagram.Message().Marshal()
		errnie.Handles(err)
		ps.SetDatagram(dg)
		return nil
	})
	defer release()

	res, err := f.Struct()
	errnie.Handles(err)
	log.Println(res.Crowd())

	return nil
}
