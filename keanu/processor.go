package keanu

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
	"github.com/theapemachine/errnie/v2"
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
	errnie.TraceIn()

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
	errnie.TraceIn()
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
	errnie.TraceIn()
	return procmap[process.Value](process.Dataframe)
}

var procmap = map[string]func(map[string]interface{}) *spdg.Datagram{
	"offline": getOffline,
}

func getOffline(df map[string]interface{}) *spdg.Datagram {
	errnie.TraceIn()

	name := viper.GetString("name")
	source := name + "." + viper.GetString(name+".stage")
	wasOffline := false

	for _, check := range lookup(&spdg.Datagram{
		Context: &spdg.Context{Annotations: []map[string]string{{"lookup": ""}}},
	}) {
		errnie.Logs.Debug(check)
	}

	return spdg.NewDatagramFromBuffer(
		"offline", source, "dgram",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"value": %v}`, wasOffline))),
	)
}

func lookup(value *spdg.Datagram) []*spdg.Datagram {
	errnie.TraceIn()

	results := make([]*spdg.Datagram, 0)

	for found := range memory.NewTree().Peek(value) {
		results = append(results, found)
	}

	errnie.TraceOut()
	return results
}
