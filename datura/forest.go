package datura

import (
	iradix "github.com/hashicorp/go-immutable-radix"
)

/*
Forest connects multiple trees together in a replication configuration
over the network, where needed.
*/
type Forest struct {
	trees []*iradix.Tree
}

/*
NewForest returns a constructed replicated in-memory radix tree that
is highly available.
*/
func NewForest(trees []*iradix.Tree) *Forest {
	return &Forest{
		trees: trees,
	}
}
