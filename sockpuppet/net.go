package sockpuppet

import (
	"net"

	"github.com/theapemachine/wrkspc/errnie"
)

func WhatIsMyIp() net.IP {
	var (
		conn net.Conn
		err  error
	)
	if conn, err = net.Dial(
		"udp", "8.8.8.8:80",
	); errnie.Handles(err) != nil {
		return nil
	}

	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}
