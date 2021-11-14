package keanu

import (
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Channel is a ingest object for a memory store. It can any type
of data (wrapped in a Datagram) and serves as a way to have data,
but not store it. Also, it's quick.
*/
type Channel struct {
	In    chan *spdg.Datagram
	scope string
	tree  *Tree
}

/*
NewChannel constructs a channel for the caller and returns a reference
to itself. Internally it holds a Graph object, which gives this channel
all the features of the in-memory index of the data lake.
*/
func NewChannel(scope string) *Channel {
	return &Channel{
		In:    make(chan *spdg.Datagram),
		scope: scope,
		tree:  NewTree(),
	}
}

/*
Cycle a channel such that it reads from the incoming Datagram provider.
The it is just a matter of storing data in the tree and we are halfway
there. The method returns itself for usabilities sake, so you can just
chain a full setup in one line: `memory.NewChannel("some scope").Cycle().In`.
*/
func (channel *Channel) Cycle(disposer twoface.Disposer) *Channel {
	go func() {
		for {
			select {
			case <-disposer.Done():
				// The Disposer gives the object that instantiated this
				// channel a way to break out of the inifnite loop even
				// though we are in a generally inaccessible goroutine.
				return
			case msg := <-channel.In:
				// We received new data on the incoming channel, poke it
				// into the radix tree.
				channel.tree.Poke(msg)
			}
		}
	}()

	return channel
}
