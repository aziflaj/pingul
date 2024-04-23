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
	lxr := &Lexer{input: []rune(input)}
	lxr.readChar()
	return lxr
}

func (l *Lexer) NextToken() token.Token {
	tkn := token.Token{Literal: l.readNextToken()}

	// handle empty literals, i.e. EOF, whitespace, newlines, tabs,	etc.
	if len(tkn.Literal) == 0 {
		if l.ch == 0 {
			tkn.Type = token.EOF
			goto getNextToken // goto statement in action. where is your god now?
		} else { // not EOF, just skip this token
			l.readChar()
			return l.NextToken()
		}
	}

	switch {
	case token.IsDelimiter(tkn.Literal):
		tkn.Type = token.Delimiters[rune(tkn.Literal[0])]
	case token.IsOperator(tkn.Literal):
		tkn.Type = token.Operators[rune(tkn.Literal[0])]
	case token.IsKeyword(tkn.Literal):
		tkn.Type = token.Keywords[string(tkn.Literal)]
	case token.IsComparisonOperator(tkn.Literal):
		tkn.Type = token.ComparisonOperators[string(tkn.Literal)]

	case token.IsInteger(tkn.Literal):
		tkn.Type = token.INT
	// TODO: Chars should be inside quotes
	// case token.IsCharacter(tkn.Literal):
	// 	tkn.Type = token.CHAR

	default:
		tkn.Type = token.IDENTIFIER
	}

getNextToken:
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

func (l *Lexer) readNextToken() []rune {
	var word []rune

	for l.ch != 0 {
		if l.ch == ' ' || l.ch == '\n' || l.ch == '\t' {
			break
		}

		// check if current character is an unary operator
		if _, ok := token.UnaryOperators[l.ch]; ok {
			// take a step back and process unary operator in the next iteration
			if len(word) > 0 {
				l.position -= 1
				l.readPosition -= 1
				break
			} else {
				word = append(word, l.ch)
				break
			}
		}

		if _, ok := token.Delimiters[l.ch]; ok {
			// take a step back and process delimiter in the next iteration
			if len(word) > 0 {
				l.position -= 1
				l.readPosition -= 1
				break
			}

			word = append(word, l.ch)
			break
		}

		word = append(word, l.ch)
		l.readChar()
	}

	return word
}
