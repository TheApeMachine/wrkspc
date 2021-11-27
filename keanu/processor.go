package keanu

import (
	"bytes"
	"fmt"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Processor is a proxy object that holds a Process it can execute and extract a value from.
Values will be respresentated as Datagram types, but what isn't these days.
*/
type Processor struct {
	process Process
}

/*
NewProcessor constructs a Processor with a Process.
*/
func NewProcessor(process, value string, df map[string]interface{}) *Processor {
	errnie.Traces()

	return &Processor{
		process: NewProcess(Synthesizer{
			Value:     value,
			Dataframe: df,
		}),
	}
}

/*
GetValue proxies the Execute method on the Process.
*/
func (processor *Processor) GetValue() *spdg.Datagram {
	errnie.Traces()
	return processor.process.Execute()
}

/*
Process is an interface that defines a common type for various Processes.
*/
type Process interface {
	Execute() *spdg.Datagram
}

/*
NewProcess constructs... And so forth.
*/
func NewProcess(processType Process) Process {
	return processType
}

/*
Synthesizer is a Process that is able to generate entirely new data points that do not natively
exist in the data landscape of a system, by inspecting various other values to infer the new data.
*/
type Synthesizer struct {
	Value     string
	Dataframe map[string]interface{}
}

/*
Execute the Process.
*/
func (process Synthesizer) Execute() *spdg.Datagram {
	errnie.Traces()
	return procmap[process.Value](process.Dataframe)
}

var procmap = map[string]func(map[string]interface{}) *spdg.Datagram{
	"offline": getOffline,
}

func getOffline(df map[string]interface{}) *spdg.Datagram {
	errnie.Traces()

	for _, check := range lookup(&spdg.Datagram{
		Context: &spdg.Context{Annotations: []spdg.Annotation{}},
	}) {
		errnie.Logs(check).With(errnie.DEBUG)
	}

	return spdg.QuickDatagram(
		spdg.DATAPOINT, "dgram",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"value": %v}`, true))),
	)
}

func lookup(value *spdg.Datagram) []*spdg.Datagram {
	errnie.Traces()

	results := make([]*spdg.Datagram, 0)

	for found := range NewTree().Peek(value) {
		results = append(results, found)
	}

	return results
}
