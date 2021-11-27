package keanu

import (
	"encoding/json"
	"strings"

	"github.com/theapemachine/errnie/v2"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Scope is an object that breaks another object into separate values
and optionally performs additional transformations, or aggregations
that are based on requirements or use-cases. This is the main way to
use the radix tree memory both read and write.
*/
type Scope struct {
	Dataframe map[string]interface{}
}

/*
NewScope returns a reference to a constructed Scope and converts a
datagram quickly to a more generic representation of key/value as strings.
This is useful mostly for the inner data, which can be of different structure.
*/
func NewScope(datagram *spdg.Datagram) *Scope {
	// Passing in nil is an actual feature, for retrieving a querying interface.
	if datagram == nil {
		return &Scope{}
	}

	df := make(map[string]interface{})
	errnie.Handles(json.Unmarshal([]byte(datagram.Data.Payload), &df))

	return &Scope{
		Dataframe: df,
	}
}

/*
Fetch is a `query` like interface to the radix tree based memory. You pass it
a `question` datagram to which in turn it will provide an `answer` datagram.

Datagram {
	Context {
		Role: "question",
		Annotations: [
			{Q: "robot/network/offline"},
			{Range: "1h"}
		]
	}
}

The above will return all the robots that were offline in the previous
hour.
*/
func (scope *Scope) Fetch(datagram *spdg.Datagram) chan *spdg.Datagram {
	return NewResult(datagram)
}

/*
Breakup is an entry method to various methods, filters, and transformations
that break up bigger datastructures into individual data points.
*/
func (scope *Scope) Breakup(keys ...string) {
	errnie.TraceIn()

	go func() {
		for _, key := range keys {
			errnie.Logs.Debug("breaking up key", key)

			proc := strings.Split(key, ".")
			value := NewProcessor(proc[0], proc[1], scope.Dataframe).GetValue()
			errnie.Logs.Debug("about to poke tree with", value)

			NewTree().Poke(value)
		}

		errnie.TraceOut()
	}()
}
