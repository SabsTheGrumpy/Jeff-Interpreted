package ast

import (
	"jeff/token"
	"testing"
)

func TestString(t *testing.T) {

	expectedString := "jeff's myVar is anotherVar;"

	program := &Program{
		Statements: []Statement{
			&JeffStatement{
				Token: token.Token{Type: token.JEFFS, Literal: "jeff's"},
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
