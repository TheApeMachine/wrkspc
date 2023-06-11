package sherlock

import "github.com/theapemachine/wrkspc/spd"

/*
Case is a container for an investigation.
*/
type Case struct {
	dg *spd.Datagram // Investigation data container.
}

/*
NewCase creates a new Case object.
*/
func NewCase(file *spd.Datagram) *Case {
	return &Case{file}
}

/*
Read implements the io.Reader interface, and every call should return
the current state of the investigation.
*/
func (c *Case) Read(p []byte) (n int, err error) {
	return 0, nil
}

/*
Write implements the io.Writer interface, and every call should append
the given data to the current state of the investigation.
*/
func (c *Case) Write(p []byte) (n int, err error) {
	return 0, nil
}

/*
Close implements the io.Closer interface, and should be called when the
investigation is complete. It should write the current state of the
investigation to the data lake.
*/
func (c *Case) Close() error {
	return nil
}
