package luthor

type Context struct {
	handler interface {
		getNextState(ContextHelper, *Lexer) LexState
	}
}

func (context Context) getNextState(lexer *Lexer) LexState {
	// errnie.Ambient().Log(errnie.DEBUG, lexer.state, "<-", lexer.token)
	newstate := context.handler.getNextState(ContextHelper{}, lexer)
	// errnie.Ambient().Log(errnie.DEBUG, newstate, "->", lexer.token)
	return newstate
}
