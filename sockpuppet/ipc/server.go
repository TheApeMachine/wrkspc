package ipc

import (
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"net"

	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

type Server struct {
	addr       string
	sock       net.Listener
	conn       net.Conn
	encryption struct {
		pubkey     *ecdsa.PublicKey
		pkbuf      []byte
		exchange   string
		algoritihm string
		cicpher    *cipher.AEAD
	}
	err error

	clients []spd.Datagram
}

func NewServer() *Server {
	return &Server{
		addr: uuid.NewString(),
	}
}

func (server *Server) Up() error {
	var conn net.Conn

	for {
		if server.conn, server.err = server.sock.Accept(); server.err != nil {
			return errnie.Handles(server.err)
		}

		var (
			priv *ecdsa.PrivateKey
		)

		if priv, server.err = ecdsa.GenerateKey(
			elliptic.P384(), rand.Reader,
		); server.err != nil {
			return errnie.Handles(server.err)
		}

		server.encryption.pubkey = &priv.PublicKey

		if !priv.IsOnCurve(
			server.encryption.pubkey.X, server.encryption.pubkey.Y,
		) {
			server.err = errors.New("keys are not on curve")
			return errnie.Handles(server.err)
		}

		server.encryption.pkbuf = elliptic.Marshal(
			elliptic.P384(),
			server.encryption.pubkey.X, server.encryption.pubkey.Y,
		)

		var (
			n   int
			buf = make([]byte, 300)
		)

		// Send the public key.
		if n, server.err = conn.Write(
			server.encryption.pkbuf,
		); server.err != nil {
			return errnie.Handles(server.err)
		}

		// Read the response from the client.
		if n, server.err = conn.Read(buf); server.err != nil {
			return errnie.Handles(server.err)
		}

		if n != 97 {
			server.err = errors.New("public key invalid length")
			return errnie.Handles(server.err)
		}

		x, y := elliptic.Unmarshal(elliptic.P384(), buf[:n])
		if key := (&ecdsa.PublicKey{elliptic.P384(), x, y}); !key.IsOnCurve(
			key.X, key.Y,
		) {
			server.err = errors.New("invalid public key")
			return errnie.Handles(server.err)
		}
	}
}

func (server *Server) exchange() [32]byte {
	var (
		shared [32]byte
		n      int
		buf    = make([]byte, 300)
	)

	// Send the public key.
	if _, server.err = server.conn.Write(
		server.encryption.pkbuf,
	); errnie.Handles(server.err) != nil {
		return shared
	}

	// Read the response from the client.
	if _, server.err = server.conn.Read(
		buf,
	); errnie.Handles(server.err) != nil {
		return shared
	}

	if n != 97 {
		server.err = errors.New("public key invalid length")
		return shared
	}

	x, y := elliptic.Unmarshal(elliptic.P384(), buf[:n])
	if key := (&ecdsa.PublicKey{elliptic.P384(), X: x, Y: y}); !key.IsOnCurve(
		key.X, key.Y,
	) {
		server.err = errors.New("invalid public key")
		return shared
	}

	return shared
}
