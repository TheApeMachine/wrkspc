package drknow

import "github.com/theapemachine/wrkspc/spd"

/*
Abstract is a nebulous object that is used as a generic, self-connected object
which is the building block for constructing the knowledge graph.
It is meant to represent a single line of reasoning or idea, as we are not trying
to model a monolithic knowledge graph, but rather a series of smaller, more focused
graphs that can be combined in a variety of ways to evaluate a mental model about the
world. The world is used as an alias for a useful scope of knowledge to achieve a goal.
*/
type Abstract struct {
	*spd.Datagram
}

/*
NewAbstract creates a new Abstract object.
*/
func NewAbstract(dg *spd.Datagram) *Abstract {
	return &Abstract{dg}
}

/*
Read implements the io.Reader interface.
*/
func (a *Abstract) Read(p []byte) (n int, err error) {
	return 0, nil
}

/*
Write implements the io.Writer interface.
*/
func (a *Abstract) Write(p []byte) (n int, err error) {
	return 0, nil
}

/*
Close implements the io.Closer interface.
*/
func (a *Abstract) Close() error {
	return nil
}
