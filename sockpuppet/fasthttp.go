package sockpuppet

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

/*
FastHTTPClient is a much faster implementation compared to the standard
library one, at the cost of not being 100% compliant.
*/
type FastHTTPClient struct {
	conn *fasthttp.Client
}

func NewFastHTTPClient() *FastHTTPClient {
	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	return &FastHTTPClient{
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
