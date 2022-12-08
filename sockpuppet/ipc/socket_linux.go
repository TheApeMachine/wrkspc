//go:build linux || darwin
// +build linux darwin

package ipc

import (
	"github.com/theapemachine/wrkspc/brazil"
)

func (server *Server) Up() error {
	addr := fmt.Sprintf("/tmp/%s.sock", server.addr)
	brazil.ClearPath(addr)

	server.listener, server.err = net.Listen("unix", addr)
}
