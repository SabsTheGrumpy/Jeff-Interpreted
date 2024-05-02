package parser

import (
	"fmt"
	"jeff/ast"
	"jeff/lexer"
	"jeff/token"
	"strconv"
)

// Order of operator precendences
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

// precedences is a mapping of token type to Precedence
var precedences = map[token.TokenType]int{
	token.EQUALS:     EQUALS,
	token.NOT_EQUALS: EQUALS,
	token.LT:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.SLASH:      PRODUCT,
	token.ASTERIX:    PRODUCT,
	token.LPAREN:     CALL,
}

type prefixParseFn func() ast.Expression

type infixParseFn func(ast.Expression) ast.Expression

// Parser is used to parse tokens into an Abstract Syntax Tree
type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates a new Parser
func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// call next token to set up current and peek token pointers (seems lazy)
	parser.nextToken()
	parser.nextToken()

	// Sets prefix parsing functions based on the token
	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefixFn(token.IDENT, parser.parseIdentifier)
	parser.registerPrefixFn(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixFn(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixFn(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixFn(token.RIGHT, parser.parseBoolean)
	parser.registerPrefixFn(token.HUANG, parser.parseBoolean)
	parser.registerPrefixFn(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefixFn(token.IF, parser.parseIfStatement)
	parser.registerPrefixFn(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefixFn(token.STRING, parser.parseStringLiteral)

	// Sets infix parsing functions based on the token
	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfixFn(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixFn(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixFn(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixFn(token.ASTERIX, parser.parseInfixExpression)
	parser.registerInfixFn(token.EQUALS, parser.parseInfixExpression)
	parser.registerInfixFn(token.NOT_EQUALS, parser.parseInfixExpression)
	parser.registerInfixFn(token.LT, parser.parseInfixExpression)
	parser.registerInfixFn(token.GT, parser.parseInfixExpression)

	// This is infix since ( in the token between ident/lit and the arguements list.
	// i.e add(2,2)
	parser.registerInfixFn(token.LPAREN, parser.parseCallExpression)

	return parser

}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Indentifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseProgram creates the AST from the input of the parser
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program

}

// parseStatement parses statements based on the keyword
// if the current token doesn't match a keyword then parse as expression
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.JEFFS:
		return p.parseJeffStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseJeffStatement() *ast.JeffStatement {

	statement := &ast.JeffStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Indentifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement

}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	statement := &ast.ExpressionStatement{Token: p.currentToken}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement

}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q to integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expresion := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.currentPrecendence()
	p.nextToken()
	expresion.Right = p.parseExpression(precedence)

	return expresion
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentToken.Type == token.RIGHT}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfStatement() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.currentToken,
	}

	// Check for (
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Move token and parse the condition
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// parse consequence
	expression.Consequence = p.parseBlockStatement()

	// if else statement exists, parse that block as well
	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression

}

// parseFunctionLiterl parses functions. e.g.
// fn (x,y) { return x + y }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currentToken}

	// Check we have the first paren for param
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Check we have the { for the body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Indentifier {
	identifiers := []*ast.Indentifier{}

	// Handle empty parameter list
	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// Grab first param
	ident := &ast.Indentifier{Token: p.currentToken, Value: p.currentToken.Literal}

	identifiers = append(identifiers, ident)

	// while there are more params, keep grabbing them
	for p.peekToken.Type == token.COMMA {
		p.nextToken() // now current token is comma
		p.nextToken()
		ident := &ast.Indentifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// Check that we have a closed paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseCallExpression parsed function calls
// these can be ident or literals.
// add(2,2) or fn(x,y) { return x + y }(2,2)
// where you have ident 'add' or literal 'fn(x, y){ return x + y }
// and then the call statement (2,2)
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {

	args := []ast.Expression{}

	// Check for empty arg list.
	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	// Parse first arg, which is an expression
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currentToken,
	}

	block.Statements = []ast.Statement{}

	p.nextToken()

	// While still in block, add statements to the statement slice
	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	} else {
		p.peekError(tokenType)
		return false
	}
}

// func (p *Parser) peekTokenIs(t token.TokenType) bool {
// 	return p.peekToken.Type == t
// }

// peekError writes error message to errors regarding tokType not matching the peekToken type
func (p *Parser) peekError(tokType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead!", tokType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function found for token %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// peekPrecedence returns the precedence of the peekToken.
// if no precedence is found for the token then default to lowest
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) currentPrecendence() int {
	if prec, ok := precedences[p.currentToken.Type]; ok {
		return prec
	}
	return LOWEST
}
