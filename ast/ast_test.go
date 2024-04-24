package ast

import (
	"jeff/token"
	"testing"
)

func TestString(t *testing.T) {

	expectedString := "let myVar = anotherVar;"

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Indentifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Indentifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != expectedString {
		t.Errorf("Incorrect program string. Expected %s but got %s", expectedString, program.String())
	}

}
