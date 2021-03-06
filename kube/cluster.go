package kube

import (
	"os/exec"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Cluster is a collection of Kubernetes Nodes.
*/
type Cluster struct {
	nodes []Node
}

/*
NewCluster constructs a new Cluster object and returns a reference to it.
*/
func NewCluster() *Cluster {
	errnie.Traces()

	return &Cluster{
		nodes: make([]Node, 0),
	}
}

/*
Up defines the steps needed to bring up the new Cluster and sends them out over a channel
so some command handler can execute them and provide feedback.
*/
func (cluster *Cluster) Up() chan *exec.Cmd {
	errnie.Traces()
	out := make(chan *exec.Cmd)
	return out
}

/*
AddNode connects a Node to this Cluster instance.
*/
func (cluster *Cluster) AddNode(node Node) *Cluster {
	errnie.Traces()

	node.Provision()
	cluster.nodes = append(cluster.nodes, node)

	return cluster
}
