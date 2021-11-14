package conquer

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/matrix"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Platform is an environment for a Command to run on.
*/
type Platform interface {
	Boot() Platform
	Parse([]string)
	Process(string) chan *spdg.Datagram
}

/*
NewPlatform constructs a Platform of the type passed in.
*/
func NewPlatform(platformType Platform) Platform {
	errnie.Traces()
	return platformType.Boot()

}

/*
Docker is a Platform for a container to run on without using Kubernetes.
*/
type Docker struct {
}

/*
Boot the runtime environment for this Platform.
*/
func (platform Docker) Boot() Platform {
	errnie.Traces()
	matrix.NewDaemon()
	return platform
}

/*
Parse the command into executable steps.
*/
func (platform Docker) Parse(command []string) {
	errnie.Traces()
	build := matrix.NewBuild(command[0])

	// Wait until the build is verified.
	<-build.Validate()
}

/*
Process the Command.
*/
func (platform Docker) Process(name string) chan *spdg.Datagram {
	errnie.Traces()
	out := make(chan *spdg.Datagram)
	return out
}
