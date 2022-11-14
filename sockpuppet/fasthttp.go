package sockpuppet

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
	"github.com/valyala/fasthttp"
)

/*
FastHTTPClient is a much faster implementation compared to the standard
library one, at the cost of not being 100% compliant.
*/
type FastHTTPClient struct {
	ctx  twoface.Context
	pool *twoface.Pool
	conn *fasthttp.Client
}

func NewFastHTTPClient() *FastHTTPClient {
	ctx := twoface.NewContext(nil)

	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	return &FastHTTPClient{
		ctx:  ctx,
		pool: twoface.NewPool(ctx).Run(),
		conn: &fasthttp.Client{
			ReadTimeout:                   readTimeout,
			WriteTimeout:                  writeTimeout,
			MaxIdleConnDuration:           maxIdleConnDuration,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:        true,
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: 5 * time.Minute,
			}).Dial,
		},
	}
}

/*
HTTPJob wraps fasthttp GET and POST requests, so we can schedule
them onto a worker pool.
*/
type HTTPJob struct {
	wg     *sync.WaitGroup
	method string
	p      []byte
}

/*
Do implements the Job interface, which enables the HTTP request to
be scheduled onto a worker pool.
*/
func (job *HTTPJob) Do() errnie.Error {
	defer job.wg.Done()

	dg := spd.Unmarshal(job.p)
	uri := spd.Payload(dg)

	url := fasthttp.AcquireURI()
	url.Parse(nil, uri)

	hc := &fasthttp.HostClient{Addr: "localhost:8080"}
	req := fasthttp.AcquireRequest()
	req.SetURI(url)
	fasthttp.ReleaseURI(url)

	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	if e := errnie.Handles(hc.Do(req, resp)); e.Type != errnie.NIL {
		return e
	}

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	job.p = resp.Body()
	return errnie.NewError(nil)
}

/*
Read implements io.Reader and in the case of this object represents
an HTTP GET request.
*/
func (client *FastHTTPClient) Read(p []byte) (n int, err error) {
	errnie.Traces()
	var wg sync.WaitGroup
	wg.Add(1)

	client.pool.Do(&HTTPJob{
		wg:     &wg,
		method: fasthttp.MethodGet,
		p:      p,
	})

	wg.Wait()
	return
}

/*
Write implements io.Writer and in the case of this object represents
an HTTP POST request.
*/
func (client *FastHTTPClient) Write(p []byte) (n int, err error) {
	errnie.Traces()
	var wg sync.WaitGroup
	wg.Add(1)

	client.pool.Do(&HTTPJob{
		wg:     &wg,
		method: fasthttp.MethodPost,
		p:      p,
	})

	wg.Wait()
	return
}

/*
Do is an incomplete, non-functional attempt to replace the default
net/http client for the S3 upload and download manager.
*/
func (client *FastHTTPClient) Do(req *http.Request) (*http.Response, error) {
	freq := fasthttp.AcquireRequest()
	freq.Header.Set("Host", "127.0.0.1:9000")
	freq.SetRequestURI(req.URL.String())
	freq.Header.SetContentLength(-1)
	freq.Header.SetMethod(req.Method)
	freq.Header.SetProtocol("HTTP/1.1")
	freq.Header.SetContentType("application/json")

	for key, value := range req.Header {
		if key != "Content-Length" {
			freq.Header.Set(key, strings.Join(value, ""))
		}
	}

	resp := fasthttp.AcquireResponse()
	err := client.conn.Do(freq, resp)
	fasthttp.ReleaseRequest(freq)

	if err == nil {
		fmt.Printf("DEBUG Response: %s\n", resp.Body())
	} else {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	defer fasthttp.ReleaseResponse(resp)

	return &http.Response{
		StatusCode:    resp.StatusCode(),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewReader(resp.Body())),
		ContentLength: int64(len(resp.Body())),
		Close:         resp.ConnectionClose(),
		Uncompressed:  false,
		Request:       req,
	}, nil
}
