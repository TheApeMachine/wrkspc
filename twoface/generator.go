package twoface

import "github.com/theapemachine/wrkspc/spdg"

/*
Generator is an interface that can be implemented by objects that somehow source and generate data.
*/
type Generator interface {
	Yield(chan *spdg.Datagram)
}

/*
NewGenerator constructs a Generator of the type that is passed in.
*/
func NewGenerator(generatorType Generator) Generator {
	return generatorType
}

/*
ParGenerator is a concurrent Generator type.
*/
type ParGenerator struct {
	Operator func(chan *spdg.Datagram)
	Disposer *Disposer
}

/*
Yield activates the generator.
*/
func (generator ParGenerator) Yield(in chan *spdg.Datagram) {
	go func() {
		for {
			select {
			case <-generator.Disposer.Done():
				// Something sent us a cancellation signal. Bail out.
				return
			default:
				// Call back to the operator function to perform another cycle of work.
				generator.Operator(in)
			}
		}
	}()
}
