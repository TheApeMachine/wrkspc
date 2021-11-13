package twoface

import "github.com/theapemachine/wrkspc/spdg"

/*
Generator is an interface that can be implemented by objects that somehow source and generate data.
*/
type Generator interface {
	Yield(chan *spdg.Datagram) chan *spdg.Datagram
}

/*
NewGenerator constructs a Generator of the type that is passed in.
*/
func NewGenerator(generatorType Generator) Generator {
	return generatorType
}

/*
GoGenerator is a concurrent Generator type.
*/
type GoGenerator struct {
	Operator func(chan *spdg.Datagram) chan *spdg.Datagram
	Disposer *Disposer
}

/*
Yield activates the generator.
*/
func (generator GoGenerator) Yield(in chan *spdg.Datagram) chan *spdg.Datagram {
	out := make(chan *spdg.Datagram)

	go func() {
		defer close(out)

		for {
			select {
			case <-generator.Disposer.Done():
				return
			default:
				out <- <-generator.Operator(in)
			}
		}
	}()

	return out
}
