package infra

import "github.com/theapemachine/wrkspc/errnie"

/*
Cluster is an interface for objects to implement that want to
orchestrate workloads.
*/
type Cluster interface {
	Provision() errnie.Error
	Teardown() errnie.Error
}

/*
NewCluster takes converts a cluster struct type into its interface
representation.
*/
func NewCluster(clusterType Cluster) Cluster {
	return clusterType
}
