package fellini

import (
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Component is an interface that objects can implement to become useable inside a Template.
*/
type Component interface {
	Initialize(chan []byte, *twoface.Disposer) Component
}

/*
NewComponent constructs a Component of the type that is passed in and calls its Initialize method.
*/
func NewComponent(
	componentType Component,
	channel chan []byte,
	disposer *twoface.Disposer,
) Component {
	return componentType.Initialize(channel, disposer)
}
