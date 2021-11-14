package conquer

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Platform is an environment for a Command to run on.
*/
type Platform interface {
	Boot() Platform
	Parse([]string) Platform
	Process() chan *spdg.Datagram
}

/*
NewPlatform constructs a Platform of the type passed in.
*/
func NewPlatform(platformType Platform) Platform {
	errnie.Traces()
	return platformType.Boot()

}
