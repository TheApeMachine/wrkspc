package spd

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure"
)

/*
Datagram acts as a wrapper around unstructured or structured data such
that it can be stored in a datalake, using the context header to keep
track of what is in the payload.
*/
type Datagram struct {
	version   string
	uuid      uuid.UUID
	checksum  string
	mime      string
	timestamp int64
	context   string
	role      string
	ptr       int
	payload   []Layer
}

/*
NewDatagram creates a new Datagram, prepares the context header,
and hands back a pointer ready for writing data.
*/
func NewDatagram(context, role, mime string) *Datagram {
	return &Datagram{
		version:   "0.0.1",
		uuid:      uuid.New(),
		timestamp: time.Now().UnixNano(),
		mime:      mime,
		context:   context,
		role:      role,
		ptr:       0,
		payload:   nil,
	}
}

/*
Read implements the io.Reader interface for Datagram, it will read
the layer at the current pointer.
*/
func (d *Datagram) Read(p []byte) (n int, err error) {
	if err != nil {
		return 0, err
	}

	return d.payload[d.ptr].Read(p)
}

/*
Write implements the io.Writer interface for Datagram, it will first
increment the pointer, then write the payload to a new layer.
*/
func (d *Datagram) Write(p []byte) (n int, err error) {
	d.ptr++
	d.payload = append(d.payload, Layer{})
	d.payload[d.ptr].Write(p)

	// Generate a checksum for the datagram.
	d.checksum, err = generateChecksum(d)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func generateChecksum(datagram *Datagram) (string, error) {
	hash, err := hashstructure.Hash(datagram, &hashstructure.HashOptions{SlicesAsSets: true})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}
