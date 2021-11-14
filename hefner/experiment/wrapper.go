package hefner

func (pipe ProtoPipe) Wrap(modeller interface{}) chan Pipe {
	if modeller != nil {
		pipe.cache.Poke(pipe.key.String(), &modeller)
	}

	return pipe.o
}

func (pipe ProtoPipe) Unwrap(modeller interface{}) chan Pipe {
	pipe.cache.Peek(pipe.key.String(), &modeller)
	return pipe.o
}

/*
Compare dynamic types in concurrent dataflow.
*/
func (pipe ProtoPipe) Compare(a chan Pipe, q interface{}) bool {
	confirm := false

	for o := range a {
		// lol. I actually can't believe Go allows me to do this, but it makes sense.
		if <-pipe.Unwrap(&q) == <-pipe.Unwrap(&o) {
			confirm = true
		}
	}

	return confirm
}
