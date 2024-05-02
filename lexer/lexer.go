package lexer

import (
	"jeff/token"
)

// Lexer. Converts characters to tokens
type Lexer struct {
	input        string
	position     int
	readPosition int
	character    byte
}

// Constructor for the lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Set the next character in the input as the current character. Then moves position pointers
// If input is at end will set character to 0
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.character = 0
	} else {
		l.character = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// newToken creates a new token of the given type and character
func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(character),
	}
}

// readIdentifier reads the next whole word from the input
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.character) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readNumber reads the next integer from the input
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.character) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// NextToken creates a token from the next set of characters (ignoring whitespace)
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.character {
	case '=':
		if l.peekChar() == '=' {
			currentChar := l.character
			l.readChar()
			t = token.Token{Type: token.EQUALS, Literal: string(currentChar) + string(l.character)}
		} 
		// else {
		// 	t = newToken(token.ASSIGN, l.character)
		// }

	case '+':
		t = newToken(token.PLUS, l.character)
	case '-':
		t = newToken(token.MINUS, l.character)
	case '!':
		if l.peekChar() == '=' {
			currentChar := l.character
			l.readChar()
			t = token.Token{Type: token.NOT_EQUALS, Literal: string(currentChar) + string(l.character)}
		} else {
			t = newToken(token.BANG, l.character)
		}
	case '*':
		t = newToken(token.ASTERIX, l.character)
	case '/':
		t = newToken(token.SLASH, l.character)
	case '<':
		t = newToken(token.LT, l.character)
	case '>':
		t = newToken(token.GT, l.character)
	case ',':
		t = newToken(token.COMMA, l.character)
	case ';':
		t = newToken(token.SEMICOLON, l.character)
	case '(':
		t = newToken(token.LPAREN, l.character)
	case ')':
		t = newToken(token.RPAREN, l.character)
	case '{':
		t = newToken(token.LBRACE, l.character)
	case '}':
		t = newToken(token.RBRACE, l.character)
	case '"':
		t.Type = token.STRING
		t.Literal = l.readString()

	case 0:
		t.Literal = ""
		t.Type = token.EOF

	default:
		if isLetter(l.character) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t
		} else if isDigit(l.character) {
			t.Literal = l.readNumber()
			t.Type = token.INT
			return t
		} else {
			t = newToken(token.ILLEGAL, l.character)
		}
	}

	l.readChar()
	return t

}

func (l *Lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\n' || l.character == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition < len(l.input) {
		return l.input[l.readPosition]
	} else {
		return 0
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.character == '"' || l.character == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(characer byte) bool {
	return 'a' <= characer && characer <= 'z' || 'A' <= characer && characer <= 'Z' || characer == '_' || characer == '\''
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
