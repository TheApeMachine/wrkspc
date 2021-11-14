package conquer

import (
	"github.com/theapemachine/wrkspc/contempt"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/kube"
	"github.com/theapemachine/wrkspc/matrix"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Kubernetes is a Platform for a container to run on without using Kubernetes.
*/
type Kubernetes struct {
	cluster *kube.Cluster
	command []string
}

/*
Boot the runtime environment for this Platform.
*/
func (platform Kubernetes) Boot() Platform {
	errnie.Traces()

	// Start containerd daemon and get a handle on our Kubernetes Cluster type.
	matrix.NewDaemon()
	platform.cluster = kube.NewCluster()

	// Start a new Scanner so we can gather the network hosts we have access to.
	scanner := contempt.NewScanner(&contempt.Range{From: 1, To: 255})

	for connection := range scanner.Sweep() {
		// Use the connection to add a new Node to the Cluster.
		platform.cluster.AddNode(
			kube.NewNode(kube.Controller{Connection: connection}),
		)
	}

	return platform
}

/*
Parse the command into executable steps.
*/
func (platform Kubernetes) Parse(command []string) Platform {
	errnie.Traces()
	platform.command = command
	return platform
}

/*
Process the Command.
*/
func (platform Kubernetes) Process() chan *spdg.Datagram {
	errnie.Traces()
	out := make(chan *spdg.Datagram)

	go func() {
		defer close(out)

	}()

	return out
}
