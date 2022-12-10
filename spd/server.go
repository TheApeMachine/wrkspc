package spd

import (
	context "context"
	"io"

	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

// RWCServer satisfies the RWC_Server interface that was generated
// by the capnp compiler.
type RWCServer struct {
	encoder  *capnp.Encoder
	decoder  *capnp.Decoder
	datagram Datagram
	err      error
}

func NewRWCServer() *RWCServer {
	errnie.Trace()
	datagram := New(APPJSN, AGGREGATION, MERGE)

	return &RWCServer{
		capnp.NewEncoder(datagram), capnp.NewDecoder(datagram),
		datagram, nil,
	}
}

func (server RWCServer) Read(ctx context.Context, call RWC_read) error {
	errnie.Trace()

	var (
		res RWC_read_Results
		buf []byte
		err error
	)

	if res, err = call.AllocResults(); err != nil {
		return errnie.Handles(err)
	}

	if buf, err = io.ReadAll(server.datagram); err != nil {
		return errnie.Handles(err)
	}

	res.SetOut(buf)

	return nil
}

// Divide is analogous to Multiply.  All capability server methods follow the
// same pattern.
func (server RWCServer) Write(ctx context.Context, call RWC_write) error {
	errnie.Trace()

	go func() {
		var (
			res RWC_write_Results
			buf []byte
			err error
		)

		if res, err = call.AllocResults(); err != nil {
			errnie.Handles(err)
			return
		}

		if buf, err = call.Args().In(); err != nil {
			errnie.Handles(err)
			return
		}

		var msg *capnp.Message
		if msg, err = server.decoder.Decode(); err != nil {
			errnie.Handles(err)
			return
		}

		var dg Datagram
		if dg, err = ReadRootDatagram(msg); err != nil {
			errnie.Handles(err)
			return
		}

		var layers capnp.DataList
		if layers, err = dg.Layers(); errnie.Handles(err) != nil {
			return
		}

	}()

	return nil
}
