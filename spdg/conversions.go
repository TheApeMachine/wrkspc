package spdg

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Encode the Datagram into a bytes.Buffer, which means it can become a Payload to another Datagram,
or can be stored in some backend storage system.
*/
func (datagram *Datagram) Encode() *bytes.Buffer {
	dgBytes := bytes.NewBuffer([]byte{})

	errnie.Handles(
		gob.NewEncoder(dgBytes).Encode(datagram),
	).With(errnie.KILL)

	return dgBytes
}

/*
Decode the Datagram from a bytes.Buffer, so we have a predictable type that holds the string
name of the type that is wrapped inside the Payload.
*/
func (datagram *Datagram) Decode() *Datagram {
	errnie.Handles(
		gob.NewDecoder(bytes.NewBuffer(datagram.Data.Payload)).Decode(datagram),
	).With(errnie.KILL)

	return datagram
}

/*
Marshal the Datagram.
*/
func (datagram *Datagram) Marshal() []byte {
	out, err := json.Marshal(datagram)
	errnie.Handles(err).With(errnie.NOOP)
	return out
}

/*
Unmarshal the Datagram.
*/
func (datagram *Datagram) Unmarshal(data []byte) *Datagram {
	errnie.Handles(json.Unmarshal(data, datagram)).With(errnie.NOOP)
	return datagram
}
