package contempt

import (
	"strconv"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Sequencer can output a sequence of something over a channel.
*/
type Sequencer struct {
	ipRange *Range
	prefix  string
}

/*
NewSequencer constructs a Sequencer.
*/
func NewSequencer(ipRange *Range, prefix string) *Sequencer {
	errnie.Traces()

	return &Sequencer{
		ipRange: ipRange,
		prefix:  prefix,
	}
}

/*
Generate the sequence and start pushing on the channel.
*/
func (seq *Sequencer) Generate() chan string {
	errnie.Traces()
	out := make(chan string)

	go func() {
		defer close(out)

		// Use the ipRange to determine the subset of the network ip addresses to sequence.
		for i := seq.ipRange.From; i < seq.ipRange.To; i++ {
			// Send out the next value in the sequence.
			out <- seq.prefix + strconv.Itoa(i)
		}
	}()

	return out
}
