package berrt

import (
	"container/ring"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Advisor defines an interface for analysis of errnie errors and should
result in the most logical answer to whether or not to continue our
processing efforts.
*/
type Advisor interface {
	Static(ring.Ring) bool
	Dynamic(<-chan errnie.Error) bool
}

/*
NewAdvisor converts a struct type to an Advisor interface type.
*/
func NewAdvisor(advisorType Advisor) Advisor {
	return advisorType
}

/*
ProtoAdvisor is the most basic and integrated version of an Advisor.
*/
type ProtoAdvisor struct{}

/*
Static analyses all errors in one cycle and bases its output on this.
*/
func (advisor ProtoAdvisor) Static(errs ring.Ring) bool {
	yc := 0
	nc := 0

	errs.Do(func(p any) {
		if p.(errnie.Error).Type != errnie.NIL {
			nc++
			return
		}

		yc++
	})

	return yc > nc
}

/*
Dynamic is a continuous process that can analyze a stream of error
values to come up with a response.
*/
func (advisor ProtoAdvisor) Dynamic(errs <-chan errnie.Error) bool {
	return true
}
