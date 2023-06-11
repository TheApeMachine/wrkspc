package luthor

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/neurosnap/sentences"
	"github.com/wrk-grp/errnie"
)

type Msg struct {
	Unstructured string `json:"unstructured"`
}

/*
Extract iterates the entire DOM and evaluates each node for any kind of text.
It should ignore code, however it should not ignore strings within code.
It should also not ignore comments is either HTML, CSS, or JS.

Each continuous string of text should be cleaned up by making sure that there
is only one space between each word, and one newline where there was one or
more newlines in the original text.
*/
func (p *Parser) Extract(origin, unstructured string) {
	errnie.Trace()

	storage := sentences.NewStorage()
	tokenizer := sentences.NewSentenceTokenizer(storage)
	sentenceTokens := tokenizer.Tokenize(unstructured)

	for _, sentenceToken := range sentenceTokens {
		errnie.Informs(sentenceToken.Text)

		// Wrap the unstructured text in a json object.
		msg := Msg{sentenceToken.Text}

		// Marshal the json object.
		buf, err := json.Marshal(msg)
		errnie.Handles(err)

		spew.Dump(buf)

		// Call the NLP service.
		res, err := p.nlpClient.Post("http://localhost:8433/analyze", "application/json", bytes.NewBuffer(buf))
		errnie.Handles(err)

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		errnie.Handles(err)

		entities := []Entity{}
		errnie.Handles(json.Unmarshal(body, &entities))

		// Add the entities to the parser.
		p.entities[origin] = append(p.entities[origin], entities...)
	}
}
