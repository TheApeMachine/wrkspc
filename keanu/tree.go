package keanu

import (
	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

var treeCache *Tree

/*
Tree is just a simple wrapper around Hashicorp's Immutable Radix Trees
that provides a simpler interface to use it in the context we use it for.
*/
type Tree struct {
	radix *iradix.Tree
}

/*
NewTree constructs a Tree if one does not exist in the cache
or refers to some existing tree in the cache which is then returned
as a reference pointer to the Tree. That means at any point you can
instantiate a Tree and you will always be talking to the same data set.
*/
func NewTree() *Tree {
	if treeCache == nil {
		treeCache = &Tree{
			radix: iradix.New(),
		}
	}

	return treeCache
}

/*
Poke the most generic of the ways to riff on a writing operation.
*/
func (tree *Tree) Poke(datagram *spdg.Datagram) {
	errnie.Logs.Debug("poking the tree with", *datagram.Prefix(), "carrying", datagram.Data.Body.Payload)

	tree.radix, _, _ = tree.radix.Insert(
		[]byte(*datagram.Prefix()),
		datagram,
	)
}

/*
Peek same same, only reading operation now. We're returning
a channel, so you probably know where this is going by now.
G e n e r a t o r s !
*/
func (tree *Tree) Peek(datagram *spdg.Datagram) chan *spdg.Datagram {
	errnie.Traces()

	if datagram == nil {
		return make(chan *spdg.Datagram)
	}

	out := make(chan *spdg.Datagram)

	go func() {
		errnie.Traces()
		defer close(out)

		// Get root node
		it := tree.radix.Root().Iterator()
		it.SeekLowerBound([]byte(
			// Search the annotations for the presence of a `lookup` key.
			twoface.Searcher(datagram.Context.Annotations, "lookup")),
		)

		// I honestly don't fully get what is going on in this for loop...
		for key, blob, ok := it.Next(); ok; key, _, ok = it.Next() {
			// Hmm, this would be even faster as a channel. Let's do that.
			out <- blob.(*spdg.Datagram).Annotate("radixKey", string(key))
		}

		errnie.Traces()
	}()

	return out
}
