package parser

import (
	"fmt"
	"strconv"

	"github.com/pmatseykanets/monkey/ast"
	"github.com/pmatseykanets/monkey/lexer"
	"github.com/pmatseykanets/monkey/token"
)

// precedence values.
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

// precedences associates token types with their precedence values.
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

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
	trace     bool
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
	p.prefixFns[token.INT] = p.parseIntegerLiteral
	p.prefixFns[token.BANG] = p.parsePrefixExpression
	p.prefixFns[token.MINUS] = p.parsePrefixExpression
	p.prefixFns[token.TRUE] = p.parseBoolean
	p.prefixFns[token.FALSE] = p.parseBoolean
	p.prefixFns[token.LPAREN] = p.parseGroupExpression

	// Register infix parsing funstions.
	p.infixFns[token.PLUS] = p.parseInfixExpression
	p.infixFns[token.MINUS] = p.parseInfixExpression
	p.infixFns[token.ASTERISK] = p.parseInfixExpression
	p.infixFns[token.SLASH] = p.parseInfixExpression
	p.infixFns[token.LT] = p.parseInfixExpression
	p.infixFns[token.GT] = p.parseInfixExpression
	p.infixFns[token.EQ] = p.parseInfixExpression
	p.infixFns[token.NOT_EQ] = p.parseInfixExpression

	// Advance twice to fill in p.curr and p.next
	p.nextToken()
	p.nextToken()

	return p
}

// WithTrace enables parsing tracing.
func (p *Parser) WithTrace() {
	p.trace = true
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
	p.errors = append(p.errors, fmt.Errorf("expected next token %s got %s", t, p.next.Type))
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
	if p.trace {
		defer untrace(trace("parseStatement"))
	}
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
	if p.trace {
		defer untrace(trace("parseLetStatement"))
	}
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
	if p.trace {
		defer untrace(trace("parseReturnStatement"))
	}
	stmt := &ast.Return{Token: p.curr}

	p.nextToken()

	// TODO: Skipping expressions until we encounter a semicolon
	for p.curr.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.BareExpr {
	if p.trace {
		defer untrace(trace("parseExpressionStatement"))
	}
	stmt := &ast.BareExpr{Token: p.curr}

	stmt.Value = p.parseExpression(LOWEST)
	if p.next.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.trace {
		defer untrace(trace("parseExpression"))
	}
	prefix := p.prefixFns[p.curr.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("missing prefixFn for %s", p.curr.Type))
		return nil
	}

	left := prefix()

	for p.next.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixFns[p.next.Type]
		if infix == nil {
			return left
		}

		p.nextToken()

		left = infix(left)
	}

	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.trace {
		defer untrace(trace("parseIdentifier"))
	}
	return &ast.Identifier{
		Token: p.curr,
		Value: p.curr.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	if p.trace {
		defer untrace(trace("parseIntegerLiteral"))
	}
	value, err := strconv.ParseInt(p.curr.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("error parsing integer literal %s", p.curr.Literal))
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.curr,
		Value: value,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.trace {
		defer untrace(trace("parsePrefixExpression"))
	}
	exp := &ast.Prefix{
		Token:    p.curr,
		Operator: p.curr.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.curr.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.next.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer untrace(trace("parseInfixExpression"))
	}
	exp := &ast.Infix{
		Token:    p.curr,
		Operator: p.curr.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	if p.trace {
		defer untrace(trace("parseBoolean"))
	}
	return &ast.Boolean{
		Token: p.curr,
		Value: p.curr.Type == token.TRUE,
	}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
