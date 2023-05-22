package luthor

import "github.com/theapemachine/wrkspc/spd"

/*
State is a data structure that indicates the current state of the
overall plan, a history of events, and any other information that
might be useful while trying to achieve a goal.
*/
type State struct {
	// PlanType indicates the current type of plan.
	PlanType PlanType
	// History is a list of events that have happened.
	History *spd.Datagram
}

/*
NewState creates a new State.
*/
func NewState(dg *spd.Datagram) *State {
	return &State{
		PlanType: SEARCHWIDE,
		History:  dg,
	}
}

/*
Read implements the io.Reader interface.
*/
func (s *State) Read(p []byte) (n int, err error) {
	return 0, nil
}

/*
Write implements the io.Writer interface.
*/
func (s *State) Write(p []byte) (n int, err error) {
	return 0, nil
}

/*
Close implements the io.Closer interface.
*/
func (s *State) Close() error {
	return nil
}
