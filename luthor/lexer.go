package luthor

import (
	"strings"

	"github.com/spf13/viper"
)

type Lexer struct {
	dataDump *DataDump
	state    LexState
	ctx      Context
	token    string
	lexeme   *Lexeme
}

func NewLexer(dataDump *DataDump) *Lexer {
	// errnie.Ambient().Log(errnie.DEBUG, "luthor.NewLexer <-", dataDump)

	if len(delimiters) == 0 {
		for _, pair := range viper.GetStringSlice("lexer.delimiters") {
			delimiters = append(delimiters, strings.Split(pair, ","))
		}
	}

	LexerTypes = viper.GetStringMapStringSlice("lexer.types")

	return &Lexer{
		dataDump: dataDump,
	}
}

func (lexer *Lexer) GenerateLexemes() chan *Lexeme {
	// errnie.Ambient().Log(errnie.DEBUG, "luthor.Lexer.GenerateLexemes <-", delimiters, LexerTypes)

	out := make(chan *Lexeme)

	go func() {
		defer close(out)

		for line := range lexer.dataDump.GenerateLines() {
			lexer.parseTokens(line, out)
		}
	}()

	return out
}

func (lexer *Lexer) parseTokens(line DataLine, out chan *Lexeme) {
	lexer.token = ""
	lexer.state = INIT
	lexer.ctx = Context{Init{}}

	for _, c := range line {
		lexer.token += string(c)

		if lexer.token == " " || lexer.token == "" {
			lexer.token = ""
			continue
		}

		lexer.state = lexer.ctx.getNextState(lexer)

		switch lexer.state {
		case INIT:
			lexer.ctx = Context{Init{}}
		case DELIMITER:
			lexer.ctx = Context{Delimiter{}}
		case LETTER:
			lexer.ctx = Context{Letter{}}
		case NUMBER:
			lexer.ctx = Context{Number{}}
		}

		if lexer.lexeme != nil {
			out <- lexer.lexeme
			lexer.lexeme = nil
		}
	}
}
