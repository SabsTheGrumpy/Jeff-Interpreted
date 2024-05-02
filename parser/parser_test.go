package parser

import (
	"fmt"
	"jeff/ast"
	"jeff/lexer"
	"testing"
)

// Test parsing jeff's statement i.e. jeff's "x = 5;"
func TestJeffStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"jeff's x is 5;", "x", 5},
		{"jeff's y is right;", "y", true},
		{"jeff's foobar is y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testJeffStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.JeffStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

// Test parsing return statements
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return right;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

// Test parsing identifiers, i.e. variable names
func TestIdentifierExpressions(t *testing.T) {

	input := "foobar"

	lexer := lexer.New(input)

	p := New(lexer)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement length of %d does not equal 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statement %T could not be converted to an ExpressionStatement", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Indentifier)

	if !ok {
		t.Fatalf("Statement Expression %T could not be converted to Identifier", statement.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("Expected Idenifier Value to be foobar, but was %s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("Expected Identifiere token literal to be foobar but was %s", ident.TokenLiteral())
	}

}

// Test parsing interger literals
func TestIntegerLiteralExpressions(t *testing.T) {
	input := "5"

	lexer := lexer.New(input)

	p := New(lexer)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Statement length of %d does not equal 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statement %T could not be converted to an ExpressionStatement", program.Statements[0])
	}

	if !testIntegerLiteral(t, statement.Expression, 5) {
		return
	}
}

func TestBoolean(t *testing.T) {
	inputTests := []struct {
		input    string
		expected bool
	}{
		{
			"huang",
			false,
		},
		{
			"right",
			true,
		},
	}

	for _, testCase := range inputTests {
		lexer := lexer.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Errorf("Expected statement length to be 1 but instead was %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statement %T could not be converted to an ExpressionStatement", program.Statements[0])
		}

		exp, ok := statement.Expression.(*ast.Boolean)

		if !ok {
			t.Fatalf("Expression %T could not be converted to boolean", statement.Expression)
		}

		if exp.Value != testCase.expected {
			t.Errorf("Expected Boolean to equal %t but was %t", testCase.expected, exp.Value)
			return
		}
	}
}

// Test parsing Prefix expressions. Like negative numbers or bang operator
func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, testCase := range prefixTests {
		lex := lexer.New(testCase.input)
		pars := New(lex)
		program := pars.ParseProgram()
		checkParserErrors(t, pars)

		if len(program.Statements) != 1 {
			t.Fatalf("Statement length of %d does not equal 1", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statement %T could not be converted to an ExpressionStatement", program.Statements[0])
		}

		prefixEx, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("Expression %T could not be converted to a PrefixExpression", statement.Expression)
		}

		if prefixEx.Operator != testCase.operator {
			t.Fatalf("Prefix Operator does not match. Expected %s but got %s", testCase.operator, prefixEx.Operator)
		}

		if !testIntegerLiteral(t, prefixEx.Right, testCase.integerValue) {
			return
		}
	}
}

// Test parsing infix expressions.
func TestParseInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, testCase := range infixTests {
		lexer := lexer.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Statement length of %d does not equal 1", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Statement %T could not be converted to an ExpressionStatement", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("Expression %T could not be converted to a InfixExpression", statement.Expression)
		}

		if !testIntegerLiteral(t, expression.Left, testCase.left) {
			return
		}

		if expression.Operator != testCase.operator {
			t.Fatalf("Expected operator to be %s, but was %s", testCase.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, testCase.right) {
			return
		}

	}
}

// Check that parser correctly sets operator precedence. i.e. that * has higher priority +
func TestOperatorPrecedence(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a + b",
			"((-a) + b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"1 * (2 + 3)",
			"(1 * (2 + 3))",
		},
	}

	for _, testCase := range tests {
		lex := lexer.New(testCase.input)
		parser := New(lex)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		actual := testCase.expected

		if actual != program.String() {
			t.Errorf("Expected %s but got %s", actual, program.String())
		}
	}

}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLiteral, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Statement Expression %T could not be converted to Identifier", il)
		return false
	}

	if intLiteral.Value != value {
		t.Fatalf("Expected Idenifier Value to be %d, but was %d", value, intLiteral.Value)
		return false
	}

	if intLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("Expected IntegerLiteral token literal to be '%d' but was %s", value, intLiteral.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {

	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser encountered %d erros", len(p.errors))
	for _, msg := range p.errors {
		t.Errorf("parser error %q", msg)
	}
	t.FailNow()
}

func testJeffStatement(t *testing.T, statement ast.Statement, name string) bool {

	if statement.TokenLiteral() != "jeff's" {
		t.Errorf("TokenLiteral not jeff's, got=%s", statement.TokenLiteral())
		return false
	}

	jeffStatement, ok := statement.(*ast.JeffStatement)
	if !ok {
		t.Errorf("Statement not JeffStatement, got=%T", statement)
		return false
	}

	if jeffStatement.Name.Value != name {
		t.Errorf("Statement name does not match. Expected %s but got %s", name, jeffStatement.Name.Value)
		return false
	}

	if jeffStatement.Name.TokenLiteral() != name {
		t.Errorf("Statement Token does not match. Expected %s but got %s", name, jeffStatement.TokenLiteral())
		return false
	}

	return true

}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Indentifier)

	if !ok {
		t.Errorf("Could not conver %T to Identifier", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("Identifier Value: %s does not match expected  %s", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Identifier Literal: %s does not match %s", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64((v)))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("Unexpected type for %T", expected)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	var valString string
	if value {
		valString = "right"
	} else {
		valString = "huang"
	}

	if bo.TokenLiteral() != valString {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("%T is not InfixExpression", exp)
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("Expected %s to equal %s", opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true

}
