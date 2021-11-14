package luthor

import (
	"unicode"
)

type FunctorType uint

const (
	CHECKER FunctorType = iota
	HANDLER
)

var stateFunctors = map[LexState]map[FunctorType]func(*Lexer, ...LexState) (LexState, bool){
	INIT:      {CHECKER: isInit, HANDLER: handleInit},
	DELIMITER: {CHECKER: isDelimiter, HANDLER: handleDelimiter},
	LETTER:    {CHECKER: isLetter, HANDLER: handleLetter},
	NUMBER:    {CHECKER: isNumber, HANDLER: handleNumber},
}

type ContextHelper struct{}

func (helper ContextHelper) CanSwitch(lexer *Lexer, lexStates ...LexState) (LexState, bool) {
	foundState := lexer.state
	ok := false

	for i, state := range lexStates {
		// Break out of the loop once we have something so we don't override with
		// a false value again.
		if ok {
			break
		}

		foundState, ok = stateFunctors[state][CHECKER](lexer, lexStates[i])
	}

	return foundState, ok
}

func isInit(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], lexer.state == INIT
}

func isDelimiter(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], inSlice(delimiters, getLastString(lexer.token, 1))
}

func isLetter(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], unicode.IsLetter(getLastChar(lexer.token, 1)[0])
}

func isNumber(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], unicode.IsNumber(getLastChar(lexer.token, 1)[0])
}

func handleInit(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], true
}

func handleDelimiter(lexer *Lexer, states ...LexState) (LexState, bool) {
	if len(lexer.token) < 2 {
		return states[0], false
	}

	isLexeme(lexer, states[0])
	return states[0], true
}

func handleLetter(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], true
}

func handleNumber(lexer *Lexer, states ...LexState) (LexState, bool) {
	return states[0], true
}
