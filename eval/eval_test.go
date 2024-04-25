package eval_test

import (
	"testing"

	"github.com/aziflaj/pingul/eval"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/object"
	"github.com/aziflaj/pingul/parser"
)

func TestEvalInt(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertIntegerObject(t, evaluated, tc.expected)
	}
}

func TestEvalBool(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestPrefixBoolNegation(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"not true", false},
		{"not not true", true},
		{"not not not true", false},

		{"not false", true},
		{"not not false", false},
		{"not not not false", true},

		// 0 is false, everything else is true
		// C is a weird language, PinguL is weirder by design!
		{"not 1", false},
		{"not 5", false},

		{"not 0", true},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestInfixBool(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"false == true", false},

		{"true == 1", true},
		{"true == 5", true},
		{"true == 0", false},

		{"false == 1", false},
		{"false == 5", false},

		{"true == not 0", true},
		{"false == not 0", false},

		{"not 0 == true", true},
		{"not 0 == false", false},

		{"1 == 2 == false", true},
		{"1 == 1 == true", true},
		{"1 < 2 == true", true},
		{"1 > 2 == false", true},
		{"1 <= 2 == true", true},
		{"1 >= 2 == true", false},
		{"not (1 >= 2) == false", false},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestInfixIntOperations(t *testing.T) {
	// intEvaledTestCases := []struct {
	// 	input    string
	// 	expected int64
	// }{
	// 	{"5 + 5", 10},
	// 	{"5 - 5", 0},
	// 	{"5 * 5", 25},
	// 	{"5 / 5", 1},
	// 	{"5 % 5", 0},
	// 	{"5 % 3", 2},

	// 	{"5 + 5 * 5", 30},
	// 	{"5 * 5 + 5", 30},
	// 	{"5 * 5 / 5", 5},

	// 	{"5 + 5 * 5 - 5", 25},
	// }

	// for _, tc := range intEvaledTestCases {
	// 	lxr := lexer.New(tc.input)
	// 	psr := parser.New(lxr)
	// 	program := psr.ParseProgram()

	// 	evaluated := eval.Eval(program)
	// 	assertIntegerObject(t, evaluated, tc.expected)
	// }

	boolEvaledTestCases := []struct {
		input    string
		expected bool
	}{
		{"5 == 5", true},
		{"5 == 1", false},
		{"5 != 5", false},
		{"5 != 1", true},

		{"nil == nil", true},
		{"nil != nil", false},

		{"5 > 5", false},
		{"5 > 1", true},
		{"5 < 5", false},
		{"5 < 1", false},
		{"5 >= 5", true},
		{"5 >= 1", true},
		{"5 <= 5", true},
		{"5 <= 1", false},
	}

	for _, tc := range boolEvaledTestCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestIfElse(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (5) { 10 }", 10},
		{"if (5 < 10) { 10 }", 10},
		{"if (5 > 10) { 10 }", nil},
		{"if (5 > 10) { 10 } else { 20 }", 20},
		{"if (5 < 10) { 10 } else { 20 }", 10},
		{"if (nil) { 10 } else {20 }", 20},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)

		integer, ok := tc.expected.(int)
		if ok {
			assertIntegerObject(t, evaluated, int64(integer))
		} else {
			_, ok := evaluated.(*object.Nil)
			if !ok {
				t.Fatalf("Expected nil object. Got=%v", evaluated)
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 1 + 2; 9;", 3},
		{"true; return 1; false;", 1},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertIntegerObject(t, evaluated, tc.expected)
	}

	input := `
if (10 > 1) {
	if (true != false) {
		return true;
	}
	return false;
}`

	lxr := lexer.New(input)
	psr := parser.New(lxr)
	program := psr.ParseProgram()

	evaluated := eval.Eval(program)
	assertBooleanObject(t, evaluated, true)
}

func TestVarStatements(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"var x = 5; x;", 5},
		{"var x = 5; var y = 10; x + y;", 15},
		{"var x = 5; var y = 10; var z = 15; x + y + z;", 30},
	}

	for _, tc := range testCases {
		lxr := lexer.New(tc.input)
		psr := parser.New(lxr)
		program := psr.ParseProgram()

		evaluated := eval.Eval(program)
		assertIntegerObject(t, evaluated, tc.expected)
	}

	// Test variable shadowing
	input := `
var x = 5;
if (true) {
	var x = 10;
	return x;
}`
	lxr := lexer.New(input)
	psr := parser.New(lxr)
	program := psr.ParseProgram()

	evaluated := eval.Eval(program)
	assertIntegerObject(t, evaluated, 10)

	input = `
var x = 5;
if (false) {
	var x = 10;
	return x;
}
return x;`
	lxr = lexer.New(input)
	psr = parser.New(lxr)
	program = psr.ParseProgram()

	evaluated = eval.Eval(program)
	assertIntegerObject(t, evaluated, 5)

	// test unasigned variable
	input = `a;`

	lxr = lexer.New(input)
	psr = parser.New(lxr)
	program = psr.ParseProgram()

	evaluated = eval.Eval(program)
	_, ok := evaluated.(*object.Nil)
	if !ok {
		t.Fatalf("Expected nil object. Got=%v", evaluated)
	}
}

///////// HELPER FUNCTIONS //////////

func assertIntegerObject(t *testing.T, obj object.Object, expected int64) {
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("Object is not an Integer. Got=%T", obj)
	}

	if integer.Value != expected {
		t.Fatalf("Object has wrong value. Got=%d, Expected=%d", integer.Value, expected)
	}
}

func assertBooleanObject(t *testing.T, obj object.Object, expected bool) {
	boolean, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("Object is not a Boolean. Got=%T", obj)
	}

	if boolean.Value != expected {
		t.Fatalf("Object has wrong value. Got=%t, Expected=%t", boolean.Value, expected)
	}
}
