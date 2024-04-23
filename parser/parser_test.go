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
		if !testLetStatement(t, stmt, tc.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
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
