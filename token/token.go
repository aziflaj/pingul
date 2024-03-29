package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = iota
	EOF

	// Identifiers + literals
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

	// Delimiters
	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Keywords
	FUNC
	VAR
)
