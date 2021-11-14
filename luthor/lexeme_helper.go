package luthor

import (
	"regexp"
	"strings"
)

func isLexeme(lexer *Lexer, state LexState) (LexState, bool) {
	var ok bool
	var n int
	var inner string

	if ok, n = isOpeningDelimiter(lexer.token); ok {
		lexer.token = string(getLastChar(lexer.token, n))
	}

	if ok, n = inDelimiter(lexer.token); !ok {
		return state, false
	}

	if inner = tokenBetweenDelimiters(lexer.token, n); inner == "" {
		return state, false
	}

	var key string

	for groupKey, tokenGroup := range LexerTypes {
		if key = searchTokenTypes(inner, groupKey, tokenGroup); key != "" {
			break
		}

		// if key = hailMaryPass(lexer.token, groupKey, tokenGroup); key != "" {
		// 	break
		// }

		key = "ANONYMOUS"
	}

	lexer.lexeme = makeLexeme(key, inner)
	lexer.token = ""

	return state, true
}

func searchTokenTypes(inner string, key string, tokenGroup []string) string {
	for _, value := range tokenGroup {
		if strings.EqualFold(value, inner) {
			return key
		}
	}

	return ""
}

func hailMaryPass(inner string, key string, tokenGroup []string) string {
	for _, value := range tokenGroup {
		pat := regexp.MustCompile(value)
		if str := pat.FindString(inner); str != "" {
			return str
		}
	}

	return ""
}

func makeLexeme(tokenType string, token string) *Lexeme {
	return &Lexeme{
		ID:    strings.ToUpper(tokenType),
		Value: strings.ToUpper(token),
	}
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
