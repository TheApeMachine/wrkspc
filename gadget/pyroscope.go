package gadget

import (
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/wrk-grp/errnie"
)

type Pyroscope struct {
	endpoint string
	profiler *pyroscope.Profiler
	err      error
}

func NewPyroscope(endpoint string) *Pyroscope {
	errnie.Trace()
	return &Pyroscope{endpoint, nil, nil}
}

func (inspector *Pyroscope) Start() error {
	errnie.Trace()
	errnie.Informs("connecting to pyroscope at", inspector.endpoint)

	// Pyroscope will give you flame charts and metrics about the
	// program resource usage. It needs to be running as a service
	// somewhere, then the endpoint should point to that service.
	inspector.profiler, inspector.err = pyroscope.Start(pyroscope.Config{
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
	})

	return inspector.err
}
