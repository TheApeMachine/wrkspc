package ford

import (
	"io"

	"github.com/theapemachine/wrkspc/twoface"
)

type Workload interface {
	io.ReadWriter
	twoface.Job
}

func NewWorkload(workloadType Workload) Workload {
	return workloadType
}
