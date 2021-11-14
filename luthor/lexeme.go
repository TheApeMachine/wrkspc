package luthor

import (
	"github.com/spf13/viper"
)

type LexState uint

const (
	INIT LexState = iota
	DELIMITER
	LETTER
	NUMBER
)

// Construct a search space based on known entities.
var LexerTypes = viper.GetStringMapStringSlice("lexer.types")

type Lexeme struct {
	ID    string
	Value interface{}
}

/*
Define pairs of delimiters between which could be interesting content.
*/
var delimiters [][]string

/*
inSlice I have been using for years. If you really want to make me
super happy, find me a way to reverse the logic so I can get my
happy path back here...
*/
func inSlice(list [][]string, item string) bool {
	for _, pairs := range list {
		for _, p := range pairs {
			if p == item {
				return true
			}
		}
	}

	return false
}
