package gadget

import (
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
)

type Pyroscope struct {
	endpoint string
}

func NewPyroscope(endpoint string) *Pyroscope {
	return &Pyroscope{endpoint}
}

func (inspector *Pyroscope) Start() {
	errnie.Informs("connecting to pyroscope at", inspector.endpoint)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: string(tweaker.GetIdentity()),

		// replace this with the address of pyroscope server
		ServerAddress: inspector.endpoint,

		// you can disable logging by setting this to nil
		Logger: nil,

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
}
