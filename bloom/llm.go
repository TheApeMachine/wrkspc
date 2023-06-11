package bloom

import (
	"bytes"
	"io"

	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/wrk-grp/errnie"
	"github.com/wrk-grp/please"
)

var models = map[string]string{
	"ner":   "dslim/bert-base-NER",
	"bloom": "bigscience/bloom",
}

type LLM struct {
	key      string
	endpoint string
	output   *bytes.Buffer
}

func NewLLM(task string) *LLM {
	errnie.Trace()

	return &LLM{
		endpoint: tweaker.GetString("bloom.endpoint") + models[task],
		key:      tweaker.GetString("bloom.key"),
		output:   bytes.NewBuffer([]byte{}),
	}
}

func (llm *LLM) Read(p []byte) (n int, err error) {
	errnie.Trace()

	if llm.output.Len() == 0 {
		return 0, io.EOF
	}

	return llm.output.Read(p)
}

func (llm *LLM) Write(p []byte) (n int, err error) {
	errnie.Trace()

	req := please.NewRequest(llm.endpoint, "POST")
	req.AddHeaders(map[string]string{
		"Authorization": "Bearer " + llm.key,
	})

	data := bytes.NewBuffer(NewMsg(string(p)).Marshal())

	io.Copy(req, data)
	errnie.Debugs("LLM: %s", llm.output.String())
	return n, nil
}

func (llm *LLM) Close() error {
	errnie.Trace()
	llm.output.Reset()
	return nil
}
