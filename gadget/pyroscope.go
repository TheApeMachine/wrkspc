package gadget

import (
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
)

type Pyroscope struct {
	endpoint string
	profiler *pyroscope.Profiler
	err      chan error
}

func NewPyroscope(endpoint string) *Pyroscope {
	errnie.Trace()
	return &Pyroscope{endpoint, nil, make(chan error)}
}

func (inspector *Pyroscope) Start() chan error {
	errnie.Trace()
	errnie.Informs("connecting to pyroscope at", inspector.endpoint)

	var err error

	// Pyroscope will give you flame charts and metrics about the
	// program resource usage. It needs to be running as a service
	// somewhere, then the endpoint should point to that service.
	if inspector.profiler, err = pyroscope.Start(pyroscope.Config{
		ApplicationName: string(tweaker.GetIdentity()),
		ServerAddress:   inspector.endpoint,
		Logger:          nil,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	}); err != nil {
		// Send the error out over the channel, so the booter will
		// unblock and allow the program to die.
		inspector.err <- err
	}

	return inspector.err
}
