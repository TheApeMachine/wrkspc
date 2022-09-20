package datura

import (
	"io"

	iradix "github.com/hashicorp/go-immutable-radix"
)

/*
Forest connects multiple trees together in a replication configuration
over the network, where needed.
*/
type Forest struct {
	trees []*iradix.Tree
	pipes []io.ReadWriter
}

/*
NewForest returns a constructed replicated in-memory radix tree that
is highly available.
*/
func NewForest(trees []*iradix.Tree, pipes []io.ReadWriter) *Forest {
	return &Forest{
		trees: trees,
		pipes: pipes,
	}
}
