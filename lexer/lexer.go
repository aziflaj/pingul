package lexer

import (
	"github.com/aziflaj/pingul/token"
)

type Lexer struct {
	input []rune

	position     int
	readPosition int

	// ch is short for char, but it's not really a char, it's a rune. sue me :)
	ch rune
}

func New(input string) *Lexer {
	lxr := &Lexer{
		input: []rune(input),
	}

	lxr.readChar()

	return lxr
}

func (l *Lexer) NextToken() token.Token {
	var tkn token.Token

	runeToTokenMap := map[rune]token.TokenType{
		'=': token.ASSIGNMENT,
		'+': token.PLUS,
		'-': token.MINUS,
		'*': token.MULTIPLY,
		'/': token.DIVIDE,
		'%': token.MODULUS,

		'(': token.LPAREN,
		')': token.RPAREN,
		'{': token.LBRACE,
		'}': token.RBRACE,

		',': token.COMMA,
		';': token.SEMICOLON,

		0: token.EOF,
	}

	if tokenType, ok := runeToTokenMap[l.ch]; ok {
		if l.ch == 0 {
			tkn = token.Token{Type: token.EOF, Literal: ""}
		} else {
			tkn = token.Token{Type: tokenType, Literal: string(l.ch)}
		}
	} else {
		tkn = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
	}

	l.readChar()
	return tkn

}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		// read the next rune
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}
