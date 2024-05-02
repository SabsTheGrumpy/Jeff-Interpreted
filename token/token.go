package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN     = "is"
	PLUS       = "+"
	MINUS      = "-"
	BANG       = "!"
	ASTERIX    = "*"
	SLASH      = "/"
	LT         = "<"
	GT         = ">"
	EQUALS     = "=="
	NOT_EQUALS = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	JEFFS    = "JEFFS"
	RIGHT    = "RIGHT"
	HUANG    = "HUANG"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"jeff's": JEFFS,
	"right":   RIGHT,
	"huang":  HUANG,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"is":     ASSIGN,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
