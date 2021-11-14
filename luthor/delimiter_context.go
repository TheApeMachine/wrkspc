package luthor

type Delimiter struct{}

func (context Delimiter) getNextState(helper ContextHelper, lexer *Lexer) LexState {
	if found, ok := helper.CanSwitch(lexer, LETTER, NUMBER); ok {
		if state, done := stateFunctors[found][HANDLER](lexer, found); done {
			return state
		}
	}

	return lexer.state
}
