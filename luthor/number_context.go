package luthor

type Number struct{}

func (context Number) getNextState(helper ContextHelper, lexer *Lexer) LexState {
	if found, ok := helper.CanSwitch(lexer, DELIMITER); ok {
		if state, done := stateFunctors[found][HANDLER](lexer, found); done {
			return state
		}
	}

	return lexer.state
}
