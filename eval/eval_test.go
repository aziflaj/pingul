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
		evaluated := evalProgram(tc.input)
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
		evaluated := evalProgram(tc.input)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestEvalString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"world"`, "world"},
		{`"hello" + " " + "world"`, "hello world"},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("Object is not a String. Got=%T", evaluated)
		}

		if string(str.Value) != tc.expected {
			t.Fatalf("Object has wrong value. Got=%s, Expected=%s", string(str.Value), tc.expected)
		}
	}
}

func TestStringConcat(t *testing.T) {
	input := `
var name = "James";
var surname = "Bond";
var fullName = name + " " + surname;
fullName;`

	evaluated := evalProgram(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not a String. Got=%T", evaluated)
	}

	expected := "James Bond"
	if string(str.Value) != expected {
		t.Fatalf("Object has wrong value. Got=%s, Expected=%s", string(str.Value), expected)
	}
}

func TestLists(t *testing.T) {
	testCases := []struct {
		input    string
		expected []int64
	}{
		{"[]", []int64{}},
		{"[1, 2, 3]", []int64{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 - 6]", []int64{3, 12, -1}},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)
		list, ok := evaluated.(*object.List)
		if !ok {
			t.Fatalf("Object is not a List. Got=%T", evaluated)
		}

		if len(list.Items) != len(tc.expected) {
			t.Fatalf("List has wrong length. Got=%d, Expected=%d", len(list.Items), len(tc.expected))
		}

		for i, item := range list.Items {
			integer, ok := item.(*object.Integer)
			if !ok {
				t.Fatalf("Object is not an Integer. Got=%T", item)
			}

			if integer.Value != tc.expected[i] {
				t.Fatalf("Object has wrong value. Got=%d, Expected=%d",
					integer.Value, tc.expected[i])
			}
		}
	}
}

func TestListIndex(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"var list = [1, 2, 3]; list[2 - 1];", 2},
		{"var list = [1, 2, 3]; var idx = 1; list[idx];", 2},
		{"var list = [1, 2, 3]; var idx = 1; list[idx] + list[idx - 1]", 3},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)
		assertIntegerObject(t, evaluated, tc.expected)
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
		evaluated := evalProgram(tc.input)
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
		evaluated := evalProgram(tc.input)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestInfixIntOperations(t *testing.T) {
	intEvaledTestCases := []struct {
		input    string
		expected int64
	}{
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"5 / 5", 1},
		{"5 % 5", 0},
		{"5 % 3", 2},

		{"5 + 5 * 5", 30},
		{"5 * 5 + 5", 30},
		{"5 * 5 / 5", 5},

		{"5 + 5 * 5 - 5", 25},
	}

	for _, tc := range intEvaledTestCases {
		evaluated := evalProgram(tc.input)
		assertIntegerObject(t, evaluated, tc.expected)
	}

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
		evaluated := evalProgram(tc.input)
		assertBooleanObject(t, evaluated, tc.expected)
	}
}

func TestIfElse(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
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
		evaluated := evalProgram(tc.input)

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
		evaluated := evalProgram(tc.input)
		assertIntegerObject(t, evaluated, tc.expected)
	}

	input := `
if (10 > 1) {
	if (true != false) {
		return true;
	}
	return false;
}`

	evaluated := evalProgram(input)
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
		{"var x = 5; if (true) { var x = 10; return x; }", 10},
		{"var x = 5; if (false) { var x = 10; return x; } return x;", 5},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)
		assertIntegerObject(t, evaluated, tc.expected)
	}

	// test unasigned variable
	input := `a;`
	evaluated := evalProgram(input)
	_, ok := evaluated.(*object.Nil)
	if !ok {
		t.Fatalf("Expected nil object. Got=%v", evaluated)
	}
}

func TestFuncs(t *testing.T) {
	input := `func(x) { x + 1; }`
	evaluated := evalProgram(input)

	fun, ok := evaluated.(*object.Func)
	if !ok {
		t.Fatalf("Expected Func object. Got=%v", evaluated)
	}

	if len(fun.Params) != 1 {
		t.Fatalf("Wrong params: %v", fun.Params)
	}

	if fun.Params[0].String() != "x" {
		t.Fatalf("Expected 'x'. Got=%s", fun.Params[0].String())
	}

	expectedBody := `{(x + 1)}`
	if fun.Body.String() != expectedBody {
		t.Fatalf("Expected %s. Got=%s", expectedBody, fun.Body.String())
	}
}

func TestFuncCalls(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"var self = func(x) { x; }; self(5);", 5},
		{"var add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"var twice = func(x) { x * 2; }; twice(5);", 10},
		{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func (x) { x % 2; }(5);", 1}, // IIFE go brr
	}

	for _, tc := range testCases {
		t.Logf("Running test case: %s", tc.input)
		evaluated := evalProgram(tc.input)
		assertIntegerObject(t, evaluated, tc.expected)
	}
}

func TestIntrinsicFuncs(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{`len("Hello, World!");`, 13},
		{`len([1, 2, 3]);`, 3},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)
		str, ok := evaluated.(*object.Integer)
		if !ok {
			t.Fatalf("Expected String object. Got=%v", evaluated)
		}

		if str.Value != tc.expected {
			t.Fatalf("Expected %d. Got=%d", tc.expected, str.Value)
		}
	}
}

func TestFibonacci(t *testing.T) {
	fibonacci := `
var fib = func(n) {
	if (n <= 1) {
		return n;
	}

	return fib(n - 1) + fib(n - 2);
};`

	testCases := []struct {
		input    string
		expected int64
	}{
		{fibonacci + "fib(0);", 0},
		{fibonacci + "fib(1);", 1},
		{fibonacci + "fib(2);", 1},
		{fibonacci + "fib(3);", 2},
		{fibonacci + "fib(4);", 3},
		{fibonacci + "fib(5);", 5},
		{fibonacci + "fib(6);", 8},
		{fibonacci + "fib(7);", 13},
		{fibonacci + "fib(8);", 21},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)
		assertIntegerObject(t, evaluated, tc.expected)
	}
}

func TestObjectLiteral(t *testing.T) {
	input := `var obj = { key: "value" }; obj`

	lxr := lexer.New(input)
	psr := parser.New(lxr)
	program := psr.ParseProgram()
	scope := object.NewScope()

	result := eval.Eval(scope, program)

	dict, ok := result.(*object.Dict)
	if !ok {
		t.Fatalf("Object is not a Dict. Got=%T", result)
	}

	if len(dict.Pairs) != 1 {
		t.Fatalf("Dict has wrong number of pairs. Got=%d, Expected=1", len(dict.Pairs))
	}

	val, ok := dict.Pairs["key"]
	if !ok {
		t.Fatalf("Dict does not have key 'key'")
	}

	strVal, ok := val.(*object.String)
	if !ok {
		t.Fatalf("Value is not a String. Got=%T", val)
	}

	if string(strVal.Value) != "value" {
		t.Fatalf("String value is wrong. Got=%s, Expected=value", string(strVal.Value))
	}
}

func TestPropertyAccess(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{`var obj = { key: "value" }; obj.key;`, "value"},
		{`var obj = { num: 42 }; obj.num;`, int64(42)},
		{`var obj = { flag: true }; obj.flag;`, true},
		{`var obj = { nested: { inner: "deep" } }; obj.nested.inner;`, "deep"},
		{`var obj = { list: [1, 2, 3] }; obj.list[1];`, int64(2)},
		{`var obj = { nested: { list: ["a", "b"] } }; obj.nested.list[0];`, "a"},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)

		switch expected := tc.expected.(type) {
		case string:
			assertStringObject(t, evaluated, expected)
		case int64:
			assertIntegerObject(t, evaluated, expected)
		case bool:
			assertBooleanObject(t, evaluated, expected)
		}
	}
}

func TestEmptyObject(t *testing.T) {
	input := `{}`

	evaluated := evalProgram(input)

	dict, ok := evaluated.(*object.Dict)
	if !ok {
		t.Fatalf("Object is not a Dict. Got=%T", evaluated)
	}

	if len(dict.Pairs) != 0 {
		t.Fatalf("Dict should be empty. Got=%d pairs", len(dict.Pairs))
	}
}

func TestObjectWithFunctions(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		// Simple function property
		{`var squarer = { perform: func(x) { x * 2 } }; squarer.perform(5);`, int64(10)},
		// Function returning string
		{`var greeter = { greet: func(name) { "Hello, " + name } }; greeter.greet("World");`, "Hello, World"},
		// Nested object with function
		{`var obj = { math: { add: func(a, b) { a + b } } }; obj.math.add(3, 4);`, int64(7)},
		// Multiple functions in object
		{`var calc = { double: func(x) { x * 2 }, square: func(x) { x * x } }; calc.double(4);`, int64(8)},
		{`var calc = { double: func(x) { x * 2 }, square: func(x) { x * x } }; calc.square(4);`, int64(16)},
		// Object with mixed properties
		{`var obj = { name: "test", fn: func(n) { n + 100 } }; obj.fn(5);`, int64(105)},
		// Accessing non-existent property returns nil
		{`var obj = { x: 1 }; obj.y;`, nil},
	}

	for _, tc := range testCases {
		evaluated := evalProgram(tc.input)

		switch expected := tc.expected.(type) {
		case string:
			assertStringObject(t, evaluated, expected)
		case int64:
			assertIntegerObject(t, evaluated, expected)
		case nil:
			_, ok := evaluated.(*object.Nil)
			if !ok {
				t.Fatalf("Expected Nil, Got=%T", evaluated)
			}
		}
	}
}

func TestObjectMethodCalls(t *testing.T) {
	input := `
		var calculator = {
			add: func(a, b) { a + b },
			multiply: func(a, b) { a * b },
			power: func(base, exp) { 
				if (exp == 0) { 
					1 
				} else { 
					base * base 
				} 
			}
		};

		calculator.add(10, 20);
	`

	evaluated := evalProgram(input)
	assertIntegerObject(t, evaluated, 30)
}

func TestComplexObjectStructure(t *testing.T) {
	input := `
		var user = {
			name: "Alice",
			age: 30,
			address: {
				street: "Main St",
				zip: 12345
			},
			getInfo: func() { "Alice is 30 years old" }
		};

		user.getInfo();
	`

	evaluated := evalProgram(input)
	assertStringObject(t, evaluated, "Alice is 30 years old")
}

func TestObjectPropertyAccess(t *testing.T) {
	input := `
		var obj = {
			data: [10, 20, 30],
			getFirst: func() { 
				obj.data[0]
			}
		};

		obj.getFirst();
	`

	evaluated := evalProgram(input)
	assertIntegerObject(t, evaluated, 10)
}

///////// HELPER FUNCTIONS //////////

func assertStringObject(t *testing.T, obj object.Object, expected string) {
	str, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("Object is not a String. Got=%T", obj)
	}

	if string(str.Value) != expected {
		t.Fatalf("String value is wrong. Got=%s, Expected=%s", string(str.Value), expected)
	}
}

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

func evalProgram(input string) object.Object {
	lxr := lexer.New(input)
	psr := parser.New(lxr)
	program := psr.ParseProgram()
	scope := object.NewScope()

	return eval.Eval(scope, program)
}
