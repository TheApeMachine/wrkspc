package luthor

type Init struct{}

func (context Init) getNextState(helper ContextHelper, lexer *Lexer) LexState {
	if found, ok := helper.CanSwitch(lexer, DELIMITER, LETTER, NUMBER); ok {
		if state, done := stateFunctors[found][HANDLER](lexer, found); done {
			return state
		}
	}

	if len(lexer.token) > 1 {
		panic("you should not be here")
	}

	return lexer.state
}
