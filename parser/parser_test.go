package parser_test

import (
	"testing"

	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/parser"
	"github.com/aziflaj/pingul/token"
)

func TestLetStatements(t *testing.T) {
	input := `var age = 28;
var name = "SpongeBob";
var result = 10 * (20 / 2);
`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()

	assertProgram(t, program)
	assertProgramLength(t, program, 3)

	testCases := []struct {
		expectedIdentifier string
	}{
		{"age"},
		{"name"},
		{"result"},
	}

	for i, tc := range testCases {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tc.expectedIdentifier) {
			return
		}
	}
}

func TestParseErrors(t *testing.T) {
	input := `var age = 28; var bob marley;`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()

	assertProgram(t, program)

	testCases := []struct {
		expectedIdentifier string
	}{
		{"age"},
	}

	for i, tc := range testCases {
		stmt := program.Statements[i]
		// assert stmt type
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Errorf("stmt not *ast.VarStatement. Got=%T", stmt)
			continue
		}

		if string(varStmt.Name.Value) != tc.expectedIdentifier {
			t.Errorf("Expected identifier %s, got %s",
				tc.expectedIdentifier, string(varStmt.Name.Value))
		}

		if len(p.Errors()) == 0 {
			t.Errorf("Expected an error, got none")
		}

		if len(p.Errors()) != 1 {
			t.Errorf("Expected 1 error, got %d", len(p.Errors()))
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
return fubar;
return if (power > 9000) { "strong" } else { "weak" };
`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()

	assertProgram(t, program)
	assertProgramLength(t, program, 3)

	for _, stmt := range program.Statements {
		// assert type
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. Got=%T", stmt)
			continue
		}

		if string(retStmt.TokenLiteral()) != "return" {
			t.Errorf("retStmt.TokenLiteral not 'return'. Got=%q",
				retStmt.TokenLiteral())
		}
	}
}

func TestStringifiedProgram(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarStatement{
				Token: token.Token{Type: token.VAR, Literal: []rune("var")},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: []rune("name")},
					Value: []rune("name"),
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: []rune("userName")},
					Value: []rune("userName"),
				},
			},
		},
	}

	expected := "var name = userName;"
	if program.String() != expected {
		t.Errorf("program.String() returned %q, expected %q",
			program.String(), expected)
	}
}

func TestParseIdentifiers(t *testing.T) {
	input := `username;`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()

	assertProgram(t, program)
	assertProgramLength(t, program, 1)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			program.Statements[0])
	}

	identExpr, ok := exprStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not *ast.Identifier. Got=%T", exprStmt.Expression)
	}

	if string(identExpr.Value) != "username" {
		t.Errorf("identExpr.Value not 'username'. Got=%v", string(identExpr.Value))
	}

	if string(identExpr.TokenLiteral()) != "username" {
		t.Errorf("identExpr.TokenLiteral not 'username'. Got=%v",
			string(identExpr.TokenLiteral()))
	}
}

func TestParseIntegerLiterals(t *testing.T) {
	input := `69420`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()

	assertProgram(t, program)
	assertProgramLength(t, program, 1)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			program.Statements[0])
	}

	intExpr, ok := exprStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression is not *ast.IntegerLiteral. Got=%T", exprStmt.Expression)
	}

	if intExpr.Value != 69420 {
		t.Errorf("intExpr.Value not '69420'. Got=%v", intExpr.Value)
	}

	if string(intExpr.TokenLiteral()) != "69420" {
		t.Errorf("intExpr.TokenLiteral not '69420'. Got=%v",
			string(intExpr.TokenLiteral()))
	}
}

// Helper functions

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if string(s.TokenLiteral()) != "var" {
		t.Errorf("s.TokenLiteral not 'var'. Got=%q", s.TokenLiteral())
		return false
	}

	// cast to the right type
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("Token not *parser.VarStatement. Got=%T", s)
		return false
	}

	if string(varStmt.Name.Value) != name {
		t.Errorf("varStmt.Name.Value not '%s'. Got=%v", name, string(varStmt.Name.Value))
		return false
	}

	if string(varStmt.Name.TokenLiteral()) != name {
		t.Errorf("varStmt.Name.Value not '%s'. Got=%v", name, string(varStmt.Name.Value))
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
}

func assertProgram(t *testing.T, program *ast.Program) {
	if program == nil {
		t.Fatalf("Program is unparsable")
	}
}

func assertProgramLength(t *testing.T, program *ast.Program, expected int) {
	if len(program.Statements) != expected {
		t.Fatalf("program.Statements does not contain %d statements. Got=%d",
			expected, len(program.Statements))
	}
}
