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

func TestPrefixString(t *testing.T) {
	tests := []struct {
		prefix *Prefix
		want   string
	}{
		{
			&Prefix{
				Token:    token.Token{Type: token.BANG, Literal: "!"},
				Operator: "!",
				Right: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "a"},
					Value: "a",
				},
			},
			"(!a)",
		},
		{
			&Prefix{
				Token:    token.Token{Type: token.MINUS, Literal: "-"},
				Operator: "-",
				Right: &Infix{
					Token: token.Token{Type: token.PLUS, Literal: "+"},
					Left: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
					Operator: "+",
					Right: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "b"},
						Value: "b",
					},
				},
			},
			"(-(5 + b))",
		},
	}

	for _, tt := range tests {
		if got := tt.prefix.String(); tt.want != got {
			t.Errorf("Expected %s got %s", tt.want, got)
		}
	}
}

func TestInfixString(t *testing.T) {
	tests := []struct {
		infix *Infix
		want  string
	}{
		{
			&Infix{
				Token: token.Token{Type: token.PLUS, Literal: "+"},
				Left: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "a"},
					Value: "a",
				},
				Operator: "+",
				Right: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "b"},
					Value: "b",
				},
			},
			"(a + b)",
		},
		{
			&Infix{
				Token: token.Token{Type: token.PLUS, Literal: "*"},
				Left: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
				Operator: "*",
				Right: &Infix{
					Token: token.Token{Type: token.PLUS, Literal: "+"},
					Left: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "a"},
						Value: "a",
					},
					Operator: "+",
					Right: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "b"},
						Value: "b",
					},
				},
			},
			"(5 * (a + b))",
		},
	}

	for _, tt := range tests {
		if got := tt.infix.String(); tt.want != got {
			t.Errorf("Expected %s got %s", tt.want, got)
		}
	}
}

func TestBooleanString(t *testing.T) {
	tests := []struct {
		boolean *Boolean
		want    string
	}{
		{
			&Boolean{
				Token: token.Token{Type: token.TRUE, Literal: "true"},
				Value: true,
			},
			"true",
		},
		{
			&Boolean{
				Token: token.Token{Type: token.FALSE, Literal: "false"},
				Value: false,
			},
			"false",
		},
	}

	for _, tt := range tests {
		if got := tt.boolean.String(); tt.want != got {
			t.Errorf("Expected %s got %s", tt.want, got)
		}
	}
}
