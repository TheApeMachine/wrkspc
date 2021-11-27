package hefner

/*
Wrap it...
*/
func (pipe ProtoPipe) Wrap(modeller interface{}) chan Pipe {
	// if modeller != nil {
	// 	pipe.cache.Poke(pipe.key.String(), &modeller)
	// }

	return pipe.o
}

/*
Your guess is as good as mine as to what this does. I wrote this a month or so back
in what could be described as a haze...
Apparently it is used somehow like this:

dg.Unwrap(&artifact) <- <-dg.Generate()

I suppose it will take a type that implements hefner Pipe and pushes its input
through the output into the Unwrap channel with a type used to model the data
into a structured type.
*/
func (pipe ProtoPipe) Unwrap(modeller interface{}) chan Pipe {
	// pipe.cache.Peek(pipe.key.String(), &modeller)
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
