package luthor

import (
	"bytes"
	"io"
	"net/http"

	"github.com/wrk-grp/errnie"
)

/*
Parser is the top level object which directs the parsing and extraction of
information from a string of HTML.
*/
type Parser struct {
	nlpClient *http.Client
	text      *bytes.Buffer
	entities  map[string][]Entity
}

/*
NewParser creates a new Parser.
*/
func NewParser() *Parser {
	errnie.Trace()

	return &Parser{
		&http.Client{},
		bytes.NewBuffer([]byte{}),
		map[string][]Entity{},
	}
}

/*
Read implements the io.Reader interface.
*/
func (p *Parser) Read(b []byte) (n int, err error) {
	errnie.Trace()
	return 0, io.EOF
}

/*
Write implements the io.Writer interface.
*/
func (p *Parser) Write(b []byte) (int, error) {
	errnie.Trace()
	return 0, nil
}

/*
Close implements the io.Closer interface.
*/
func (p *Parser) Close() error {
	errnie.Trace()
	return nil
}
