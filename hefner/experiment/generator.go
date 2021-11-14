package hefner

/*
Generator orchestrates data flow with unknown types.
*/
type Generator interface {
	Generate() chan Pipe
	Wrap(interface{}) chan Pipe
	Unwrap(interface{}) chan Pipe
	Compare(chan Pipe, interface{}) bool
}

/*
Generate some data flow by piping the input to the output.
*/
func (pipe ProtoPipe) Generate() chan Pipe {
	go func() {
		defer close(pipe.o)

		for i := range pipe.i.Generate() {
			pipe.o <- i
		}
	}()

	return pipe.o
}
