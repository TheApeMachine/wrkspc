package spdg

/*
Header combines meta data from HTTP Headers, internal definitions, and other sources
into a generalistic descriptor referencing teh Payload of a Datagram.
*/
type Header struct {
	Sources []string
}

/*
NewHeader constructs the Header for a Datagram.
*/
func NewHeader() *Header {
	return &Header{}
}
