package ast

import (
	"testing"

	"github.com/pmatseykanets/monkey/token"
)

func TestProgramString(t *testing.T) {
	prg := &Program{
		Statements: []Statement{
			&Let{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if want, got := `let myVar = anotherVar;`, prg.String(); want != got {
		t.Errorf("Expected %s got %s", want, got)
	}
}
