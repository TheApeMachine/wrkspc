package errnie

/*
Write implements the io.Writer interface.
*/
func (ctx *Context) Write(p []byte) (n int, err error) {
	return ctx.output.Write(p)
}

/*
Close implements the io.Closer interface.
*/
func (ctx *Context) Close() error {
	return Handles(ctx.fh.Close())
}
