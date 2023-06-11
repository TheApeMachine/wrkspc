package bloom

import (
	"encoding/json"

	"github.com/wrk-grp/errnie"
)

type Parameters struct {
	AggregationStrategy string `json:"aggregation_strategy"`
}

type Msg struct {
	Inputs     string     `json:"inputs"`
	Parameters Parameters `json:"parameters"`
}

func NewMsg(input string) *Msg {
	errnie.Trace()

	return &Msg{
		Inputs: input,
		Parameters: Parameters{
			AggregationStrategy: "average",
		},
	}
}

func (msg *Msg) Marshal() []byte {
	errnie.Trace()

	buf, err := json.Marshal(msg)
	errnie.Handles(err)
	return buf
}

type Result struct {
	GeneratedText string `json:"generated_text"`
}
