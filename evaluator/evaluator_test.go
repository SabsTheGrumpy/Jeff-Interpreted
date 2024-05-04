package evaluator

import (
	"jeff/lexer"
	"jeff/object"
	"jeff/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{
			"5", 5,
		},
		{
			"10", 10,
		},
		{
			"-5",
			-5,
		},
		{
			"-10",
			-10,
		},
		{
			"10 + 5",
			15,
		},
		{
			"10 - 5",
			5,
		},
		{
			"10 * 5",
			50,
		},
		{
			"10 / 5",
			2,
		},
		{
			"10 / 5 + 2",
			4,
		},
		{
			"10 / 5 + 2 * 5",
			12,
		},
	}

	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		testIntegerObject(t,evaluated,testCase.expected)
	}
}


func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{
			"right",
			true,
		},
		{
			"huang",
			false,
		},
		{
			"5 > 1",
			true,
		},
		{
			"5 < 1",
			false,
		},
		{
			"1 > 5",
			false,
		},
		{
			"1 < 5",
			true,
		},
		{
			"5 == 1",
			false,
		},
		{
			"5 != 1",
			true,
		},
		{
			"5 == 5",
			true,
		},
		{
			"5 != 5",
			false,
		},
		{
			"right == right",
			true,
		},
		{
			"right == huang",
			false,
		},
		{
			"right != right",
			false,
		},
		{
			"right != huang",
			true,
		},
		{
			"huang == huang",
			true,
		},
		{
			"(1 < 4) == right",
			true,
		},
		{
			"(1 > 4) == right",
			false,
		},
		{
			"(1 < 4) == huang",
			false,
		},
		{
			"(1 > 4) == huang",
			true,
		},

		

	}


	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		testBooleanObject(t, evaluated, testCase.expected)
	}
}


func TestBangOperator(t *testing.T) {

	tests := []struct {
		input string
		expected bool
	}{
		{
			"!right",
			false,
		},
		{
			"!huang",
			true,
		},
		{
			"!5",
			false,
		},
		{
			"!!5",
			true,
		},
	}

	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		testBooleanObject(t, evaluated, testCase.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{"if (right) { 10 }", 10},
		{"if (huang) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if( 1 > 2) {10} else {20}", 20},
		{"if( 2 > 1) {10} else {20}", 10},
	}

	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		integer, ok := testCase.expected.(int)
		if !ok {
			testNullObject(t, evaluated)
		} else {
			testIntegerObject(t, evaluated, int64(integer))
		}
	}
}

func TestReturnStatements(t *testing.T) {

	tests := []struct {
		input string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 10;", 10},
		{"if(10 > 1) { if(10 > 1) { return 10;} return 1}", 10},
	}

	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		testIntegerObject(t, evaluated, testCase.expected)
	}

}


func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		expectedMessage string
	}{
		{"5 + right", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + right; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-right", "unknown operator: -BOOLEAN"},
		{"right + huang", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; right + huang; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}



	for _, testCase := range tests {
		evaluated := testEval(testCase.input)

		errObj, ok := evaluated.(*object.ERROR)
		if !ok {
			t.Errorf("no error object returned. Got %T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != testCase.expectedMessage {
			t.Errorf("Unexpected message. Expected %s but got %s", testCase.expectedMessage, errObj.Message)
		}


	}
}

func TestJeffStatements(t *testing.T) {

	tests := []struct{
		input string
		expected int64
	}{
		{"jeff's x is 5; x", 5},
		{"jeff's a is 5 * 5;a;", 25},
		{"jeff's a is 5; jeff's b is 5; jeff's c is a * b; c;", 25},
	}

	for _, testCase := range tests {
		testIntegerObject(t, testEval(testCase.input), testCase.expected)
	}

}


func TestFunctionObject(t *testing.T) {
	input := "fn(x) {  x + 2 ; };"


	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("Evaluated %T(%+v) could not be converted to function", evaluated, evaluated)
	}


	if len(fn.Parameters) != 1 {
		t.Fatalf("Expected lengh of parameters to be 1, but got %d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("Expected function parameter to be x but got %s", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Expected function body to be %s but got %s", expectedBody, fn.Body)
	}
}


func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"jeff's identity is fn(x) { x; }; identity(5);", 5},
		{"jeff's identity is fn(x) { return x; }; identity(5);", 5},
		{"jeff's double is fn(x) { x * 2; }; double(5);", 10},
		{"jeff's add is fn(x, y) { x + y; }; add(5, 5);", 10},
		{"jeff's add is fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, testCase := range tests {
		testIntegerObject(t, testEval(testCase.input), testCase.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}


func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.ERROR)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}





func testNullObject(t *testing.T, obj object.Object) bool {

	if obj != NULL {
		t.Errorf("Object %T (%+v) is not null", obj, obj)
		return false
	}
	return true

}


func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok  {
		t.Errorf("object %T (%+v) is not boolean", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. Expected %t, but got %t", expected, result.Value)
		return false
	}

	return true
}


func testEval(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok  {
		t.Errorf("object %T (%+v) is not integer", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. Expected %d, but got %d", expected, result.Value)
		return false
	}

	return true
}