package kube

import (
	"github.com/theapemachine/wrkspc/errnie"
)

/*
ClusterType defines an enum which select the type of cluster we will be
interfacing with.
*/
type ClusterType uint

const (
	KIND ClusterType = iota
	REMOTE
)

/*
Cluster is a wrapper around a Kubernetes cluster.
There are various types of cluster that could be referenced by this wrapper,
including KIND (Kubernetes In Docker) for local setups, and remote clusters
when interfacing with staging/QA/Production environments.
*/
type Cluster struct {
	IsProvisioned bool
}

/*
NewCluster constructs, or connects to, a Kubernetes environment so we have
a location where we can deploy our containers.
*/
func NewCluster() *Cluster {
	return &Cluster{}
}

/*
Provision the selected cluster so we can deploy containers onto it.
*/
func (cluster *Cluster) Provision() *errnie.Error {
	// Start a new KIND (Kubernetes In Docker) cluster for a local setup.
	// I have implemented the kind/pkg/log/types interfaces in errnie, so we
	// can keep a consistent terminal log experience.
	/*
		return errnie.NewError(app.Run(
			errnie.GetErrnie(), cmd.StandardIOStreams(), []string{
				"create",
				"cluster",
				"--name", "wrkspc",
				"--config", "cmd/cfg/.kind-config.yml",
			},
		))
	*/
	return nil
}

/*
Teardown brings everything back down.
*/
func (cluster *Cluster) Teardown() *errnie.Error {
	/*
		return errnie.NewError(app.Run(
			errnie.GetErrnie(), cmd.StandardIOStreams(), []string{
				"delete",
				"cluster",
				"--name", "wrkspc",
			},
		))
	*/
	return nil
}
