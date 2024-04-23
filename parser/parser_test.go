package parser_test

import (
	"strconv"
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

func TestParsePrefixExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"not 5;", "not", 5},
		{"-15;", "-", 15},
		{"not true", "not", true},
		{"not false", "not", false},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		p := parser.New(lxr)
		program := p.ParseProgram()

		assertProgram(t, program)
		checkParserErrors(t, p)
		assertProgramLength(t, program, 1)

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
				program.Statements[0])
		}

		prefixExpr, ok := exprStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression is not *ast.PrefixExpression. Got=%T", exprStmt.Expression)
		}

		if string(prefixExpr.Operator) != tc.operator {
			t.Errorf("prefixExpr.Operator not '%s'. Got=%v", tc.operator, string(prefixExpr.Operator))
		}

		if !testLiteralExpression(t, prefixExpr.Right, tc.value) {
			return
		}
	}
}

func TestParseInfixExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"true == true", true, "==", true},
	}

	for _, tc := range testCases {
		lex := lexer.New(tc.input)
		p := parser.New(lex)
		program := p.ParseProgram()

		assertProgram(t, program)
		checkParserErrors(t, p)
		assertProgramLength(t, program, 1)

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
				program.Statements[0])
		}

		testInfixExpression(t, exprStmt.Expression, tc.left, tc.operator, tc.right)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-(a)) * b)"},
		{"a * -b", "(a * (-(b)))"},
		{"not -a", "(not((-(a))))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-(5)) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"3 * (2 - 6)", "(3 * (2 - 6))"},
		{"b / (a + c)", "(b / (a + c))"},
		{"not (a + b)", "(not((a + b)))"},
		{"not a + b", "((not(a)) + b)"},
		{"a + not b", "(a + (not(b)))"},
		{"not (true == false)", "(not((true == false)))"},
	}

	for index, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		t.Logf("Test %d: %s", index, tt.input)
		checkParserErrors(t, p)
		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpressions(t *testing.T) {
	input := `if (x > y) { true } else { false }`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()
	assertProgram(t, program)
	checkParserErrors(t, p)
	assertProgramLength(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			program.Statements[0])
	}

	ifExpr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. Got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", ">", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		for _, s := range ifExpr.Consequence.Statements {
			t.Logf("Statement: %s", s.String())
		}
		t.Errorf("Consequence does not have 1 statement. Got=%d",
			len(ifExpr.Consequence.Statements))
	}

	consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			ifExpr.Consequence.Statements[0])
	}

	if !testBoolean(t, consequence.Expression, true) {
		return
	}

	if len(ifExpr.Alternative.Statements) != 1 {
		t.Errorf("Alternative does not have 1 statement. Got=%d",
			len(ifExpr.Alternative.Statements))
	}

	alternative, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			ifExpr.Alternative.Statements[0])
	}

	if !testBoolean(t, alternative.Expression, false) {
		return
	}
}

func TestFuncExpressions(t *testing.T) {
	input := `func(x, y) { x + y; }`

	lxr := lexer.New(input)
	p := parser.New(lxr)
	program := p.ParseProgram()
	assertProgram(t, program)
	checkParserErrors(t, p)
	assertProgramLength(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			program.Statements[0])
	}

	funcExpr, ok := stmt.Expression.(*ast.FuncExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FuncExpression. Got=%T", stmt.Expression)
	}

	if len(funcExpr.Params) != 2 {
		t.Errorf("funcExpr.Params does not have 2 elements. Got=%d", len(funcExpr.Params))
	}

	testCases := []struct {
		expectedParam string
	}{
		{"x"},
		{"y"},
	}

	for i, tc := range testCases {
		if string(funcExpr.Params[i].Value) != tc.expectedParam {
			t.Errorf("Expected param %s, got %s",
				tc.expectedParam, string(funcExpr.Params[i].Value))
		}
	}

	if len(funcExpr.Body.Statements) != 1 {
		t.Errorf("funcExpr.Body does not have 1 statement. Got=%d",
			len(funcExpr.Body.Statements))
	}

	bodyStmt, ok := funcExpr.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("funcExpr.Body.Statements[0] is not *ast.ExpressionStatement. Got=%T",
			funcExpr.Body.Statements[0])
	}

	if !testInfixExpression(t, bodyStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestFuncParams(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"func() {}", []string{}},
		{"func(x) {}", []string{"x"}},
		{"func(x, y) {}", []string{"x", "y"}},
		{"func(x, y, z) {}", []string{"x", "y", "z"}},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		p := parser.New(lxr)
		program := p.ParseProgram()

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got=%T",
				program.Statements[0])
		}

		funcExpr, ok := stmt.Expression.(*ast.FuncExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.FuncExpression. Got=%T", stmt.Expression)
		}

		if len(funcExpr.Params) != len(tc.expected) {
			t.Errorf("funcExpr.Params does not have %d elements. Got=%d",
				len(tc.expected), len(funcExpr.Params))
		}

		for i, param := range funcExpr.Params {
			if string(param.Value) != string(tc.expected[i]) {
				t.Errorf("Expected param %s, got %s",
					string(tc.expected[i]), string(param.Value))
			}
		}
	}
}

///////// Helper functions /////////

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

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBoolean(t, expr, v)
	}

	t.Errorf("type of expr not handled. Got=%T", expr)
	return false
}

func testBoolean(t *testing.T, expr ast.Expression, value bool) bool {
	boolExpr, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr not *ast.Boolean. Got=%T", expr)
		return false
	}

	if boolExpr.Value != value {
		t.Errorf("boolExpr.Value not %t. Got=%t", value, boolExpr.Value)
		return false
	}

	if string(boolExpr.TokenLiteral()) != strconv.FormatBool(value) {
		t.Errorf("boolExpr.TokenLiteral not %t. Got=%v", value, string(boolExpr.TokenLiteral()))
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intExpr, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. Got=%T", il)
		return false
	}

	if intExpr.Value != value {
		t.Errorf("intExpr.Value not %d. Got=%d", value, intExpr.Value)
		return false
	}

	if string(intExpr.TokenLiteral()) != strconv.Itoa(int(value)) {
		t.Errorf("intExpr.TokenLiteral not %d. Got=%v", value, string(intExpr.TokenLiteral()))
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	identExpr, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr not *ast.Identifier. Got=%T", expr)
		return false
	}

	if string(identExpr.Value) != value {
		t.Errorf("identExpr.Value not %s. Got=%v", value, string(identExpr.Value))
		return false
	}

	if string(identExpr.TokenLiteral()) != value {
		t.Errorf("identExpr.TokenLiteral not %s. Got=%v",
			value, string(identExpr.TokenLiteral()))
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	expr ast.Expression,
	left interface{}, operator string, right interface{},
) bool {
	infixExpr, ok := expr.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("expression is not *ast.InfixExpression. Got=%T", expr)
		return false
	}

	if !testLiteralExpression(t, infixExpr.Left, left) {
		return false
	}

	if string(infixExpr.Operator) != operator {
		t.Errorf("infixExpr.Operator not '%s'. Got=%v", operator, string(infixExpr.Operator))
		return false
	}

	if !testLiteralExpression(t, infixExpr.Right, right) {
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

	for _, stmt := range program.Statements {
		t.Logf("Statement: %s", stmt.String())
	}

}
