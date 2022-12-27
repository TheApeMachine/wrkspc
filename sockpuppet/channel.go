package sockpuppet

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

type ChannelType uint

const (
	ERROR ChannelType = iota
)

type Channel struct {
	ctx *twoface.Context
	rwc io.ReadWriteCloser
	t   ChannelType
	err error
	out chan error
}

func NewChannel(
	ctx *twoface.Context, rwc io.ReadWriteCloser, t ChannelType,
) *Channel {
	return &Channel{ctx, rwc, t, nil, make(chan error)}
}

func (channel *Channel) Read(p []byte) (n int, err error) {
	if n, err = channel.rwc.Read(p); err != nil && err != io.EOF {
		channel.err = err
		channel.Close()
	}

	return n, io.EOF
}

func (channel *Channel) Write(p []byte) (n int, err error) {
	if n, err = channel.rwc.Write(p); err != nil && err != io.EOF {
		channel.err = err
		channel.Close()
	}

	return n, io.EOF
}

func (channel *Channel) Close() error {
	defer close(channel.out)
	errnie.Handles(channel.rwc.Close())
	channel.out <- channel.err
	return errnie.Handles(channel.err)
}
