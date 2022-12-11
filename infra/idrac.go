package infra

import (
	"bytes"
	"net"

	"github.com/theapemachine/wrkspc/errnie"
	"golang.org/x/crypto/ssh"
)

/*
IDRAC is a struct that represents a remote management controller
for a Dell server.
*/
type IDRAC struct {
	ip       string
	username string
	password string
}

/*
NewIDRAC returns a pointer to an instance of IDRAC.
*/
func NewIDRAC(ip, username, password string) *IDRAC {
	return &IDRAC{ip, username, password}
}

/*
Reboot the server remotely using its IDRAC management interface.

As long as you are on the same network, or the IDRAC interface is
exposed to the internet and you have the correct IP address, this
allows you to reboot the system, even when it is turned off.
*/
func (idrac *IDRAC) Reboot() {
	var (
		client  *ssh.Client
		session *ssh.Session
		err     error
		out     = bytes.NewBuffer([]byte{})
		buf     []byte
	)

	// Connect to the IDRAC interface over a Secure Shell Connection.
	if client, err = ssh.Dial("tcp", idrac.ip+":22", &ssh.ClientConfig{
		User: idrac.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(idrac.password),
		},
		HostKeyCallback: func(host string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}); errnie.Handles(err) != nil {
		return
	}

	// Verify the SSH session.
	if session, err = client.NewSession(); errnie.Handles(err) != nil {
		return
	}

	defer session.Close()

	session.Stdout = out

	// Reboot the machine via the IDRAC remote management controller.
	if buf, err = session.Output("racadm serveraction hardreset"); err != nil {
		return
	}

	out.Write(buf)
}
