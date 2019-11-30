package lexer

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pmatseykanets/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x,y) {
	x + y;
};

let result = add(five, ten);

!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	tests := []struct {
		want token.Token
	}{
		{token.Token{Type: token.LET, Literal: "let"}},
		{token.Token{Type: token.IDENT, Literal: "five"}},
		{token.Token{Type: token.ASSIGN, Literal: "="}},
		{token.Token{Type: token.INT, Literal: "5"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.LET, Literal: "let"}},
		{token.Token{Type: token.IDENT, Literal: "ten"}},
		{token.Token{Type: token.ASSIGN, Literal: "="}},
		{token.Token{Type: token.INT, Literal: "10"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.LET, Literal: "let"}},
		{token.Token{Type: token.IDENT, Literal: "add"}},
		{token.Token{Type: token.ASSIGN, Literal: "="}},
		{token.Token{Type: token.FUNCTION, Literal: "fn"}},
		{token.Token{Type: token.LPAREN, Literal: "("}},
		{token.Token{Type: token.IDENT, Literal: "x"}},
		{token.Token{Type: token.COMMA, Literal: ","}},
		{token.Token{Type: token.IDENT, Literal: "y"}},
		{token.Token{Type: token.RPAREN, Literal: ")"}},
		{token.Token{Type: token.LBRACE, Literal: "{"}},
		{token.Token{Type: token.IDENT, Literal: "x"}},
		{token.Token{Type: token.PLUS, Literal: "+"}},
		{token.Token{Type: token.IDENT, Literal: "y"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.RBRACE, Literal: "}"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.LET, Literal: "let"}},
		{token.Token{Type: token.IDENT, Literal: "result"}},
		{token.Token{Type: token.ASSIGN, Literal: "="}},
		{token.Token{Type: token.IDENT, Literal: "add"}},
		{token.Token{Type: token.LPAREN, Literal: "("}},
		{token.Token{Type: token.IDENT, Literal: "five"}},
		{token.Token{Type: token.COMMA, Literal: ","}},
		{token.Token{Type: token.IDENT, Literal: "ten"}},
		{token.Token{Type: token.RPAREN, Literal: ")"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.BANG, Literal: "!"}},
		{token.Token{Type: token.MINUS, Literal: "-"}},
		{token.Token{Type: token.SLASH, Literal: "/"}},
		{token.Token{Type: token.ASTERISK, Literal: "*"}},
		{token.Token{Type: token.INT, Literal: "5"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.INT, Literal: "5"}},
		{token.Token{Type: token.LT, Literal: "<"}},
		{token.Token{Type: token.INT, Literal: "10"}},
		{token.Token{Type: token.GT, Literal: ">"}},
		{token.Token{Type: token.INT, Literal: "5"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.IF, Literal: "if"}},
		{token.Token{Type: token.LPAREN, Literal: "("}},
		{token.Token{Type: token.INT, Literal: "5"}},
		{token.Token{Type: token.LT, Literal: "<"}},
		{token.Token{Type: token.INT, Literal: "10"}},
		{token.Token{Type: token.RPAREN, Literal: ")"}},
		{token.Token{Type: token.LBRACE, Literal: "{"}},
		{token.Token{Type: token.RETURN, Literal: "return"}},
		{token.Token{Type: token.TRUE, Literal: "true"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.RBRACE, Literal: "}"}},
		{token.Token{Type: token.ELSE, Literal: "else"}},
		{token.Token{Type: token.LBRACE, Literal: "{"}},
		{token.Token{Type: token.RETURN, Literal: "return"}},
		{token.Token{Type: token.FALSE, Literal: "false"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.RBRACE, Literal: "}"}},
		// 10 == 10;
		// 10 != 9;
		{token.Token{Type: token.INT, Literal: "10"}},
		{token.Token{Type: token.EQ, Literal: "=="}},
		{token.Token{Type: token.INT, Literal: "10"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.INT, Literal: "10"}},
		{token.Token{Type: token.NOT_EQ, Literal: "!="}},
		{token.Token{Type: token.INT, Literal: "9"}},
		{token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.Token{Type: token.EOF, Literal: ""}},
	}

	l := New(strings.NewReader(input))

	for i, tt := range tests {
		got := l.NextToken()
		if !cmp.Equal(tt.want, got) {
			t.Fatalf("[Test %d] Expected %v got %v", i, tt.want, got)
		}
	}
}

func TestLexerPeek(t *testing.T) {
	input := "10 == 10;"

	l := New(strings.NewReader(input))
	l.readNext()
	r := l.peek()
	if want, got := '0', r; want != got {
		t.Errorf("Expected peek value %v got %v", want, got)
	}
	l.readNext()
	if want, got := '0', l.r; want != got {
		t.Errorf("Expected read value %v got %v", want, got)
	}
}
