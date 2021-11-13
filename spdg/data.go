package spdg

import "bytes"

/*
Data is the wrapper around the inner data and type being transported by the Datagram.
The Header will hold meta data regarding the Payload that holds the raw bytes of the inner data
inside of a bytes.Buffer type.
*/
type Data struct {
	Header  *Header
	Payload []byte
}

/*
NewData constructs a Data object to carry the Header and Payload.
*/
func NewData(header *Header, payload *bytes.Buffer) *Data {
	return &Data{Header: header, Payload: payload.Bytes()}
}
