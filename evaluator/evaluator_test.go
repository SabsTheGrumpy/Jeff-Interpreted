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
			"true",
			true,
		},
		{
			"false",
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
			"true == true",
			true,
		},
		{
			"true == false",
			false,
		},
		{
			"true != true",
			false,
		},
		{
			"true != false",
			true,
		},
		{
			"false == false",
			true,
		},
		{
			"(1 < 4) == true",
			true,
		},
		{
			"(1 > 4) == true",
			false,
		},
		{
			"(1 < 4) == false",
			false,
		},
		{
			"(1 > 4) == false",
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
			"!true",
			false,
		},
		{
			"!false",
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

	return Eval(program)
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