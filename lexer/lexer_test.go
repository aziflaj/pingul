package lexer_test

import (
	"testing"

	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGNMENT, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lxr := lexer.New(input)

	for i, testToken := range tests {
		token := lxr.NextToken()

		if token.Type != testToken.expectedType {
			t.Fatalf("tests[%d] - wrong token type. Expected=%v, got=%v", i, testToken.expectedType, token.Type)
		}

		if token.Literal != testToken.expectedLiteral {
			t.Fatalf("tests[%d] - wrong token literal. Expected=%v, got=%v", i, testToken.expectedLiteral, token.Literal)
		}
	}
}
