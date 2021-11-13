package kube

import (
	"github.com/theapemachine/wrkspc/contempt"
	"github.com/theapemachine/wrkspc/errnie"
	"sigs.k8s.io/kind/pkg/cluster"
)

/*
Node is an interface objects can implement to describe a Kube Node.
*/
type Node interface {
	Initialize() Node
}

/*
NewNode constructs a node of the type that is passed in.
*/
func NewNode(nodeType Node) Node {
	errnie.Traces()
	return nodeType.Initialize()
}

/*
Controller is a Kube Node type that provides a control plane for Workers.
*/
type Controller struct {
	Connection contempt.Connection
}

/*
Initialize the Node.
*/
func (node Controller) Initialize() Node {
	errnie.Traces()

	provider := cluster.NewProvider()

	// Build a Kind Cluster running inside a Container.
	errnie.Logs("building local Kubernetes cluster").With(errnie.INFO)
	errnie.Handles(provider.Create("kind")).With(errnie.KILL)
	errnie.Logs("cluster has gone up").With(errnie.INFO)

	return node
}

/*
Worker is a Kube Node type that only does Worker.
*/
type Worker struct {
	Connection contempt.Connection
}

/*
Initialize the Node.
*/
func (node Worker) Initialize() Node {
	errnie.Traces()
	return node
}
