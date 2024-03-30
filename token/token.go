package token

import "fmt"

type TokenType int

type Token struct {
	Type    TokenType
	Literal []rune
}

const (
	ILLEGAL = iota
	EOF

	// Identifiers & Literals
	IDENTIFIER
	INT
	CHAR

	// Operators
	ASSIGNMENT
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULUS

	// Comparison
	EQUAL
	NOT_EQUAL
	GREATER_THAN
	LESS_THAN
	GREATER_THAN_OR_EQUAL
	LESS_THAN_OR_EQUAL

	// Logical

	// Delimiters
	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Keywords
	NIL
	VAR
	FUNC
	RETURN
	TRUE
	FALSE
	AND
	OR
	NOT
	IF
	ELSE
)

var Keywords = map[string]TokenType{
	"nil":    NIL,
	"var":    VAR,
	"func":   FUNC,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"and":    AND,
	"or":     OR,
	"not":    NOT,
	"if":     IF,
	"else":   ELSE,
}

var Delimiters = map[rune]TokenType{
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	',': COMMA,
	';': SEMICOLON,
}

var Operators = map[rune]TokenType{
	'=': ASSIGNMENT,
	'+': PLUS,
	'-': MINUS,
	'*': MULTIPLY,
	'/': DIVIDE,
	'%': MODULUS,
}

var ComparisonOperators = map[string]TokenType{
	"==": EQUAL,
	"!=": NOT_EQUAL,
	">":  GREATER_THAN,
	"<":  LESS_THAN,
	">=": GREATER_THAN_OR_EQUAL,
	"<=": LESS_THAN_OR_EQUAL,
}

func IsComparisonOperator(r []rune) bool {
	_, ok := ComparisonOperators[string(r)]
	return ok
}

func IsInteger(literal []rune) bool {
	if len(literal) == 0 {
		return false
	}

	for _, char := range string(literal) {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// idk if i even need this but it's here now
func IsCharacter(literal []rune) bool {
	return len(literal) == 1
}

func IsOperator(r []rune) bool {
	if len(r) != 1 {
		return false
	}
	_, ok := Operators[r[0]]
	return ok
}

func IsDelimiter(r []rune) bool {
	if len(r) != 1 {
		return false
	}
	_, ok := Delimiters[r[0]]
	return ok
}

func IsKeyword(literal []rune) bool {
	_, ok := Keywords[string(literal)]
	return ok
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v, '%v')", t.Type, string(t.Literal))
}
