package main

import (
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/theapemachine/wrkspc/cmd"
	"github.com/theapemachine/wrkspc/errnie"
)

func main() {
	// Start a connection to a Pyroscope server and
	// collect metrics on the performance of wrkspc.
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "theapemachine.wrkspc.app",
		ServerAddress:   "http://localhost:4040",
		Logger:          nil,

		ProfileTypes: []pyroscope.ProfileType{
			// These profile types are enabled by default.
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// These profile types are optional.
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

	// Entrypoint to the CLI handling.
	errnie.Kills(cmd.Execute())
}
