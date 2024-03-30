package lexer_test

import (
	"testing"

	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 28;
var timeGoesBy = func(currentAge, yearsPassed) {
	return currentAge + yearsPassed;
};

var newAge = timeGoesBy(age, 1);
age < newAge;
age <= newAge;
age > newAge;
age >= newAge;
age == newAge;
age != newAge;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral []rune
	}{
		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("age")},
		{token.ASSIGNMENT, []rune("=")},
		{token.INT, []rune("28")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("timeGoesBy")},
		{token.ASSIGNMENT, []rune("=")},
		{token.FUNC, []rune("func")},
		{token.LPAREN, []rune("(")},
		{token.IDENTIFIER, []rune("currentAge")},
		{token.COMMA, []rune(",")},
		{token.IDENTIFIER, []rune("yearsPassed")},
		{token.RPAREN, []rune(")")},
		{token.LBRACE, []rune("{")},
		{token.RETURN, []rune("return")},
		{token.IDENTIFIER, []rune("currentAge")},
		{token.PLUS, []rune("+")},
		{token.IDENTIFIER, []rune("yearsPassed")},
		{token.SEMICOLON, []rune(";")},
		{token.RBRACE, []rune("}")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.ASSIGNMENT, []rune("=")},
		{token.IDENTIFIER, []rune("timeGoesBy")},
		{token.LPAREN, []rune("(")},
		{token.IDENTIFIER, []rune("age")},
		{token.COMMA, []rune(",")},
		{token.INT, []rune("1")},
		{token.RPAREN, []rune(")")},
		{token.SEMICOLON, []rune(";")},

		{token.IDENTIFIER, []rune("age")},
		{token.LESS_THAN, []rune("<")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.SEMICOLON, []rune(";")},

		{token.IDENTIFIER, []rune("age")},
		{token.LESS_THAN_OR_EQUAL, []rune("<=")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.SEMICOLON, []rune(";")},

		{token.IDENTIFIER, []rune("age")},
		{token.GREATER_THAN, []rune(">")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.SEMICOLON, []rune(";")},

		{token.IDENTIFIER, []rune("age")},
		{token.GREATER_THAN_OR_EQUAL, []rune(">=")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.SEMICOLON, []rune(";")},

		{token.IDENTIFIER, []rune("age")},
		{token.EQUAL, []rune("==")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.SEMICOLON, []rune(";")},

		{token.IDENTIFIER, []rune("age")},
		{token.NOT_EQUAL, []rune("!=")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.SEMICOLON, []rune(";")},

		{token.EOF, []rune("")},
	}

	lxr := lexer.New(input)

	for i, testToken := range tests {
		token := lxr.NextToken()

		if token.Type != testToken.expectedType {
			t.Fatalf("tests[%d] - wrong token type. Expected=%v, got=%v", i, testToken.expectedType, token.Type)
		}

		if string(token.Literal) != string(testToken.expectedLiteral) {
			t.Fatalf("tests[%d] - wrong token literal. Expected=%v, got=%v", i, testToken.expectedLiteral, token.Literal)
		}
	}
}
