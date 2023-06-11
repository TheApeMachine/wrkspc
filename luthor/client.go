package luthor

import (
	"bytes"
	"io"

	"github.com/wrk-grp/errnie"
	"github.com/wrk-grp/please"
)

/*
NLPClient is a wrapper around a please Request that calls the NLP service.
*/
type NLPClient struct {
	*please.Request
	buffer *bytes.Buffer
}

/*
NewNLPClient creates a new NLPClient.
*/
func NewNLPClient() *NLPClient {
	errnie.Trace()

	return &NLPClient{
		please.NewRequest(
			"http://localhost:8433",
			"POST",
		),
		bytes.NewBuffer([]byte{}),
	}
}

/*
Read implements the io.Reader interface.
*/
func (nlp *NLPClient) Read(b []byte) (n int, err error) {
	errnie.Trace()

	if nlp.buffer.Len() > 0 {
		return nlp.buffer.Read(b)
	}

	return 0, io.EOF
}

/*
Write implements the io.Writer interface.
*/
func (nlp *NLPClient) Write(b []byte) (int, error) {
	errnie.Trace()
	return nlp.buffer.Write(b)
}
