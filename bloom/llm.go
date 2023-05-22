package bloom

import (
	"encoding/json"
	"fmt"

	"github.com/theapemachine/am/network"
	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/wrk-grp/errnie"
)

type LLM struct {
	req *network.Request
}

func NewLLM() *LLM {
	errnie.Trace()

	endpoint := tweaker.GetString("models.bloom.endpoint")
	key := tweaker.GetString("models.bloom.key")
	req := network.NewRequest(network.POST, endpoint)
	req.AddHeader("Authorization", "Bearer "+key)
	req.AddHeader("Content-Type", "application/json")

	return &LLM{req}
}

func (llm *LLM) Predict(input []map[string]string) chan string {
	errnie.Trace()

	out := make(chan string)

	go func() {
		defer close(out)

		res := []Result{}
		msgs := ""

		for _, in := range input {
			for key, val := range in {
				msgs += fmt.Sprintf("%s: %s\n", key, val)
			}
		}

		msg := llm.req.Do(NewMsg(msgs).Marshal())

		errnie.Handles(json.Unmarshal(
			msg,
			&res,
		))

		out <- res[0].GeneratedText
	}()

	return out
}