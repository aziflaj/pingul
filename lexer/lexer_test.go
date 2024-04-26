package lexer_test

import (
	"testing"

	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

func TestNextToken(t *testing.T) {
	input := `var fullName = "Mr. Squarepants, Spongebob";
var age = 28;
var timeGoesBy = func(currentAge, yearsPassed) {
	return currentAge + yearsPassed;
};

var newAge = timeGoesBy(age, 1);
if (age < newAge) {
	return (1 == 1);
} else {
	return (1 != 1);
}

var truthness = (age <= newAge) and (2 >= 1);
var falseness = (age > newAge) or (2 < 1);

var amIAlive = true and not false;

var x = 1;
var y = -x;
var z = x + y;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral []rune
	}{
		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("fullName")},
		{token.ASSIGNMENT, []rune("=")},
		{token.STRING, []rune("Mr. Squarepants, Spongebob")},
		{token.SEMICOLON, []rune(";")},

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

		{token.IF, []rune("if")},
		{token.LPAREN, []rune("(")},
		{token.IDENTIFIER, []rune("age")},
		{token.LESS_THAN, []rune("<")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.RPAREN, []rune(")")},
		{token.LBRACE, []rune("{")},
		{token.RETURN, []rune("return")},
		{token.LPAREN, []rune("(")},
		{token.INT, []rune("1")},
		{token.EQUAL, []rune("==")},
		{token.INT, []rune("1")},
		{token.RPAREN, []rune(")")},
		{token.SEMICOLON, []rune(";")},
		{token.RBRACE, []rune("}")},
		{token.ELSE, []rune("else")},
		{token.LBRACE, []rune("{")},
		{token.RETURN, []rune("return")},
		{token.LPAREN, []rune("(")},
		{token.INT, []rune("1")},
		{token.NOT_EQUAL, []rune("!=")},
		{token.INT, []rune("1")},
		{token.RPAREN, []rune(")")},
		{token.SEMICOLON, []rune(";")},
		{token.RBRACE, []rune("}")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("truthness")},
		{token.ASSIGNMENT, []rune("=")},
		{token.LPAREN, []rune("(")},
		{token.IDENTIFIER, []rune("age")},
		{token.LESS_THAN_OR_EQUAL, []rune("<=")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.RPAREN, []rune(")")},
		{token.AND, []rune("and")},
		{token.LPAREN, []rune("(")},
		{token.INT, []rune("2")},
		{token.GREATER_THAN_OR_EQUAL, []rune(">=")},
		{token.INT, []rune("1")},
		{token.RPAREN, []rune(")")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("falseness")},
		{token.ASSIGNMENT, []rune("=")},
		{token.LPAREN, []rune("(")},
		{token.IDENTIFIER, []rune("age")},
		{token.GREATER_THAN, []rune(">")},
		{token.IDENTIFIER, []rune("newAge")},
		{token.RPAREN, []rune(")")},
		{token.OR, []rune("or")},
		{token.LPAREN, []rune("(")},
		{token.INT, []rune("2")},
		{token.LESS_THAN, []rune("<")},
		{token.INT, []rune("1")},
		{token.RPAREN, []rune(")")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("amIAlive")},
		{token.ASSIGNMENT, []rune("=")},
		{token.TRUE, []rune("true")},
		{token.AND, []rune("and")},
		{token.NOT, []rune("not")},
		{token.FALSE, []rune("false")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("x")},
		{token.ASSIGNMENT, []rune("=")},
		{token.INT, []rune("1")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("y")},
		{token.ASSIGNMENT, []rune("=")},
		{token.MINUS, []rune("-")},
		{token.IDENTIFIER, []rune("x")},
		{token.SEMICOLON, []rune(";")},

		{token.VAR, []rune("var")},
		{token.IDENTIFIER, []rune("z")},
		{token.ASSIGNMENT, []rune("=")},
		{token.IDENTIFIER, []rune("x")},
		{token.PLUS, []rune("+")},
		{token.IDENTIFIER, []rune("y")},
		{token.SEMICOLON, []rune(";")},

		{token.EOF, []rune("")},
	}

	lxr := lexer.New(input)

	for i, tt := range tests {
		tkn := lxr.NextToken()

		testToken := token.Token{
			Type:    tt.expectedType,
			Literal: tt.expectedLiteral,
		}

		if tkn.Type != testToken.Type {
			t.Fatalf("tests[%d] - wrong Token Type. Expected=%v, got=%v",
				i, testToken, tkn)
		}

		if string(tkn.Literal) != string(testToken.Literal) {
			t.Fatalf("tests[%d] - wrong Token Literal. Expected=%v, got=%v",
				i, testToken, tkn)
		}
	}
}
