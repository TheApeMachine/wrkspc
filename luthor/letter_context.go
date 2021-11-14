package luthor

type Letter struct{}

func (context Letter) getNextState(helper ContextHelper, lexer *Lexer) LexState {
	if found, ok := helper.CanSwitch(lexer, DELIMITER, NUMBER); ok {
		if state, done := stateFunctors[found][HANDLER](lexer, found); done {
			return state
		}
	}

	return lexer.state
}
