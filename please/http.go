package please

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
HTTP is a Request transported over the HTTP protocol.
*/
type HTTP struct {
	client    *http.Client
	transport *http.Transport
	host      string
	path      string
}

/*
Do the Request.
*/
func (request HTTP) Do(options *spdg.Datagram) chan *spdg.Datagram {
	request = request.ensureClient()
	out := make(chan *spdg.Datagram)

	go func() {
		resp, err := request.client.Get(request.host + request.path)
		errnie.Handles(err).With(errnie.NOOP)

		body, err := ioutil.ReadAll(resp.Body)
		errnie.Handles(err).With(errnie.NOOP)

	out <- spdg.QuickDatagram(spdg.ANONYMOUS, "json", bytes.NewBuffer(body))
	}()

	return out
}

/*
ensureClient will check and possibly configure the HTTP client built into Go.
*/
func (request HTTP) ensureClient() HTTP {
	if request.transport == nil {
		request.transport = &http.Transport{}
	}

	if request.client == nil {
		request.client = &http.Client{Transport: request.transport}
	}

	return request
}
