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
