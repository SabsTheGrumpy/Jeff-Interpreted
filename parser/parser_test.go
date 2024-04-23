package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)


// Test parsing let statement i.e. let "x = 5;"
func TestLetStatement(t *testing.T) {

	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	lex := lexer.New(input)
	pars := New(lex)

	program := pars.ParseProgram()
	checkParserErrors(t, pars)

	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements but got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, testCase := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, testCase.expectedIdentifier) {
			return
		}

	}

}


// Test parsing return statements
func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 92828;
	`

	lex := lexer.New(input)
	pars := New(lex)

	program := pars.ParseProgram()
	checkParserErrors(t, pars)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements but got %d", len(program.Statements))
	}


	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Statement not ReturnStatement. Got=%T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("ReturnStatement token not 'return'. Got=%s", returnStatement.TokenLiteral())
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

	statement, ok  := program.Statements[0].(*ast.ExpressionStatement)

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

	statement, ok  := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statement %T could not be converted to an ExpressionStatement", program.Statements[0])
	}

	if !testIntegerLiteral(t, statement.Expression, 5) {
		return
	}
}


// Test parsing Prefix expressions. Like negative numbers or bang operator
func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
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
	
		statement, ok  := program.Statements[0].(*ast.ExpressionStatement)
	
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
		input string
		left int64
		operator string
		right int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5 , "*", 5},
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
		input string
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


func testLetStatement(t *testing.T, statement ast.Statement, name string) bool  {
		
	if statement.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral not Let, got=%s", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement not LetStatement, got=%T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("Statement name does not match. Expected %s but got %s", name, letStatement.Name.Value)
		return false
	}


	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("Statement Token does not match. Expected %s but got %s", name, letStatement.TokenLiteral())
		return false
	}


	return true

}


func testIdentifier(t *testing.T, exp ast.Expression, value string) bool{
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
	}
	t.Errorf("Unexpected type for %T", expected)
	return false
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

