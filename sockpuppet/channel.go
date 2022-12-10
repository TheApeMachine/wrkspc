package sockpuppet

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

type ChannelType uint

const (
	ERROR ChannelType = iota
)

type Channel struct {
	rwc io.ReadWriteCloser
	t   ChannelType
}

func NewChannel(rwc io.ReadWriteCloser, t ChannelType) *Channel {
	return &Channel{rwc, t}
}

func (channel *Channel) Read(p []byte) (n int, err error) {
	return
}

func (channel *Channel) Write(p []byte) (n int, err error) {
	return
}

func (channel *Channel) Close() error {
	return errnie.Handles(nil)
}
