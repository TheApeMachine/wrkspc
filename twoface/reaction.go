package twoface

/*
Reaction is a type that specializes in conditional logic.
*/
type Reaction interface {
	Cause(...bool) chan bool
	Effect() chan bool
}

/*
NewReaction constructs a Reaction of the type that is passed in.
*/
func NewReaction(reactionType Reaction) Reaction {
	return reactionType
}

/*
Boolean looks for and acts upon any true or false values.
The idea is to add it to a type in the constructor and reuse it.

It eliminates the need to use the `if` keyword.
It does its work concurrently and reports over a channel.

USAGE:
  replaces:
    if someCondition == 1 && !anotherCondition && ... {
		doSomething()
	}
  with:
    reaction.Effect(effectorFn, <- reaction.Cause(
		someCondition == 1, !anotherCondition, ...,
	))
*/
type Boolean struct {
	Disposer *Disposer
}

/*
Cause returns any value that was false.
*/
func (reaction Boolean) Cause(values ...bool) chan bool {
	out := make(chan bool)

	go func() {
		defer close(out)
		defer reaction.Disposer.Cleanup()

		var ok bool

		// Receive a value and a reduced by one slice of bool.
		for ok, values = reaction.iterate(values); !ok; {
			// value was false, send to returner.
			out <- ok
		}
	}()

	return out
}

/*
Effect produces a bahavior when it receives a false value.
*/
func (reaction Boolean) Effect(effector func(bool), next chan bool) {
	go func() {
		for {
			select {
			case cause := <-next:
				effector(cause)
			case <-reaction.Disposer.Done():
				return
			}
		}
	}()
}

/*
iterate pops a value of the slice of bool and returns it,
plus the reduced with one slice of bool.
*/
func (reaction Boolean) iterate(values []bool) (bool, []bool) {
	return values[0], values[1:]
}
