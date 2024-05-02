package lexer

import (
	"jeff/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `jeff's five is 5;
jeff's ten is 10;

jeff's add is fn(x, y) {
	x + y;
};
!-/*5;
5 < 10 > 5;

jeff's result is add(five, ten);
right huang if else return

if(x == 3)
if(x != 3)

`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.JEFFS, "jeff's"},
		{token.IDENT, "five"},
		{token.ASSIGN, "is"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.JEFFS, "jeff's"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "is"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.JEFFS, "jeff's"},
		{token.IDENT, "add"},
		{token.ASSIGN, "is"},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERIX, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.JEFFS, "jeff's"},
		{token.IDENT, "result"},
		{token.ASSIGN, "is"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RIGHT, "right"},
		{token.HUANG, "huang"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.EQUALS, "=="},
		{token.INT, "3"},
		{token.RPAREN, ")"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.NOT_EQUALS, "!="},
		{token.INT, "3"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, testCase := range tests {
		tok := lexer.NextToken()
		if tok.Type != testCase.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, testCase.expectedType, tok.Type)
		}

		if tok.Literal != testCase.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, testCase.expectedLiteral, tok.Literal)
		}
	}
}
