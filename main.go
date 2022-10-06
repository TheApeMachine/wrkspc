package main

import (
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/theapemachine/wrkspc/cmd"
	"github.com/theapemachine/wrkspc/errnie"
)

func main() {
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "theapemachine.wrkspc.app",
		ServerAddress:   "http://localhost:4040",
		Logger:          nil,

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

	errnie.Handles(cmd.Execute())
}
