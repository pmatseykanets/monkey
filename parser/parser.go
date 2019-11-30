package parser

import (
	"fmt"
	"github.com/pmatseykanets/monkey/ast"
	"github.com/pmatseykanets/monkey/lexer"
	"github.com/pmatseykanets/monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x  or !x
	CALL        // foo(x)

)

type prefixFn func() ast.Expression
type infixFn func(ast.Expression) ast.Expression

// Parser .
type Parser struct {
	lex       *lexer.Lexer
	curr      token.Token
	next      token.Token
	errors    []error
	prefixFns map[token.TokenType]prefixFn
	infixFns  map[token.TokenType]infixFn
}

// New creates a new instance of Parser.
func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		lex:       lex,
		errors:    make([]error, 0),
		prefixFns: make(map[token.TokenType]prefixFn),
		infixFns:  make(map[token.TokenType]infixFn),
	}

	// Register prefix parsing funstions.
	p.prefixFns[token.IDENT] = p.parseIdentifier

	// Advance twice to fill in p.curr and p.next
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curr = p.next
	p.next = p.lex.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.next.Type != t {
		p.peekError(t)
		return false
	}

	p.nextToken()
	return true
}

// Errors returns a slice of errors.
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Errorf("Expected next token %s got %s", t, p.next.Type))
}

// Parse .
func (p *Parser) Parse() *ast.Program {
	prg := &ast.Program{}
	prg.Statements = []ast.Statement{}

	for p.curr.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prg.Statements = append(prg.Statements, stmt)
		}
		p.nextToken()
	}

	return prg
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curr.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.Let {
	stmt := &ast.Let{Token: p.curr}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curr, Value: p.curr.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Skipping expressions until we encounter a semicolon
	for p.curr.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.Return {
	stmt := &ast.Return{Token: p.curr}

	p.nextToken()

	// TODO: Skipping expressions until we encounter a semicolon
	for p.curr.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.BareExpr {
	stmt := &ast.BareExpr{Token: p.curr}

	stmt.Value = p.parseExpression(LOWEST)
	if p.next.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFns[p.curr.Type]
	if prefix == nil {
		return nil
	}

	left := prefix()

	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curr,
		Value: p.curr.Literal,
	}
}
