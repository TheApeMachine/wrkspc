package passepartout

import "io"

/*
Store is anything that implements io.ReadWriter.
*/
type Store interface {
	io.ReadWriter
}
