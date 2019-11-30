package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/pmatseykanets/monkey/token"
)

// Lexer breaks text read from input into a stream of tokens.
type Lexer struct {
	input *bufio.Reader
	pos   int
	r     rune
	err   error
}

// New creates a new instance of Lexer.
func New(input io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(input)}
	return l
}

// FromString is a named constructor that creates a lexer from a string.
func FromString(s string) *Lexer {
	return New(strings.NewReader(s))
}

func (l *Lexer) readNext() {
	r, sz, err := l.input.ReadRune()
	l.err = err
	l.r = r
	l.pos += sz
}

func (l *Lexer) peek() rune {
	r, _, err := l.input.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0
		}
		l.err = err
		return 0
	}
	err = l.input.UnreadRune()
	if err != nil {
		l.err = err
		return 0
	}
	return r
}

func (l *Lexer) Error() error {
	return l.err
}

// NextToken consumes and returns the next token.
func (l *Lexer) NextToken() token.Token {
	if l.pos == 0 {
		l.readNext()
	}
	l.skipWhitespace()
	if l.err == io.EOF {
		return token.Token{Type: token.EOF, Literal: ""}
	}

	tok := token.Token{Literal: string(l.r)}
	switch l.r {
	case '=':
		if l.peek() == '=' {
			l.readNext()
			tok.Literal = "=="
			tok.Type = token.EQ
			break
		}
		tok.Type = token.ASSIGN
	case ';':
		tok.Type = token.SEMICOLON
	case '(':
		tok.Type = token.LPAREN
	case ')':
		tok.Type = token.RPAREN
	case ',':
		tok.Type = token.COMMA
	case '+':
		tok.Type = token.PLUS
	case '-':
		tok.Type = token.MINUS
	case '/':
		tok.Type = token.SLASH
	case '*':
		tok.Type = token.ASTERISK
	case '!':
		if l.peek() == '=' {
			l.readNext()
			tok.Literal = "!="
			tok.Type = token.NOT_EQ
			break
		}
		tok.Type = token.BANG
	case '<':
		tok.Type = token.LT
	case '>':
		tok.Type = token.GT
	case '{':
		tok.Type = token.LBRACE
	case '}':
		tok.Type = token.RBRACE
	default:
		if isLetter(l.r) {
			tok.Literal = l.readIdent()
			tok.Type = token.IdentType(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.r) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		tok.Type = token.ILLEGAL
	}

	l.readNext()
	return tok
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (l *Lexer) readIdent() string {
	var s string
	for isLetter(l.r) {
		s += string(l.r)
		l.readNext()
		if l.err != nil {
			break
		}
	}
	return s
}

func (l *Lexer) skipWhitespace() {
	if l.err != nil {
		return
	}

	for unicode.IsSpace(l.r) {
		l.readNext()
		if l.err != nil {
			return
		}
	}
}

func (l *Lexer) readNumber() string {
	var s string
	for unicode.IsDigit(l.r) {
		s += string(l.r)
		l.readNext()
		if l.err != nil {
			break
		}
	}
	return s
}
