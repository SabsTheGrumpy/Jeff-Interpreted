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
	}

	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		testIntegerObject(t,evaluated,testCase.expected)
	}

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