package drknow

import (
	"fmt"
	"io"

	"github.com/arriqaaq/art"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

type Tree struct {
	root *art.Tree
}

func NewTree() *Tree {
	return &Tree{root: art.NewTree()}
}

func (tree *Tree) Read(p []byte) (n int, err error) {
	tree.root.Scan(p, func(n *art.Node) {
		if n.IsLeaf() {
			errnie.Debugs(string(n.Value().([]byte)))
		}
	})

	return len(p), err
}

func (tree *Tree) Write(p []byte) (n int, err error) {
	var dg *spd.Datagram
	if err = dg.Decode(p); err != nil {
		return n, errnie.Handles(err)
	}

	if !tree.root.Insert(dg.Prefix().Bytes(), p) {
		return n, errnie.Handles(fmt.Errorf(
			"failed to insert %v into tree at %s",
			p, dg.Prefix().String(),
		))
	}

	return len(p), io.EOF
}

func (tree *Tree) Close() error {
	return errnie.Handles(nil)
}
