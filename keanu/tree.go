package keanu

import (
	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

var treeCache *Tree

/*
Tree wraps hashicorp's immutable radix tree to provide what is essentially a key/value store.
*/
type Tree struct {
	radix *iradix.Tree
}

/*
NewTree constructs a Tree and stores it in a singleton cache such that subsequent calls to NewTree
will always return the same object. Should there be a need to segment Tree based stores, the
suggestion would be to use a map where the key is your segmentation type and the value (or values)
are the (pointer to the) Tree.
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
Poke the Tree with some data wrapped in a Datagram.
*/
func (tree *Tree) Poke(datagram *spdg.Datagram) {
	errnie.Traces()

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
