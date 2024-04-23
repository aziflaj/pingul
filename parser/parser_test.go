package parser_test

import (
	"testing"

	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/parser"
)

func TestLetStatements(t *testing.T) {
	input := `var age = 28;
var name = "SpongeBob";
var result = 10 * (20 / 2);
`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got=%d",
			len(program.Statements))
	}

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

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. Got=%d",
			len(program.Statements))
	}

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

return func() {
	return 5;
};
`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got=%d",
			len(program.Statements))
	}

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
