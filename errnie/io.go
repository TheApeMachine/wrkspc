package errnie

import "io"

/*
Read implements the io.Reader interface.
*/
func (ctx *Context) Read(p []byte) (n int, err error) {
	n = copy(p, ctx.log.Bytes())
	ctx.log.Reset()
	return n, io.EOF
}

/*
Write implements the io.Writer interface.
*/
func (ctx *Context) write(p []byte) (n int, err error) {
	return ctx.log.Write(p)
}

/*
Close implements the io.Closer interface.
*/
func (ctx *Context) Close() error {
	defer func() {
		ctx.log = nil
		ctx = nil
	}()

	return nil
}
