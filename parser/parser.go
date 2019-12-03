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
	token.LPAREN:   CALL,
}

type prefixFn func() ast.Expression
type infixFn func(ast.Expression) ast.Expression

// Parser parses a stream of tokens produced by lexer
// and returns an AST of the program.
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
	p.prefixFns[token.IF] = p.parseIfExpression
	p.prefixFns[token.FUNCTION] = p.parseFunctionLiteral

	// Register infix parsing funstions.
	p.infixFns[token.PLUS] = p.parseInfixExpression
	p.infixFns[token.MINUS] = p.parseInfixExpression
	p.infixFns[token.ASTERISK] = p.parseInfixExpression
	p.infixFns[token.SLASH] = p.parseInfixExpression
	p.infixFns[token.LT] = p.parseInfixExpression
	p.infixFns[token.GT] = p.parseInfixExpression
	p.infixFns[token.EQ] = p.parseInfixExpression
	p.infixFns[token.NOT_EQ] = p.parseInfixExpression
	p.infixFns[token.LPAREN] = p.parseCallExpression

	// Advance twice to fill in p.curr and p.next.
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

func (p *Parser) expectNext(t token.TokenType) bool {
	if p.next.Type != t {
		p.errors = append(p.errors, fmt.Errorf("expected token type %s got %s", t, p.next.Type))
		return false
	}

	p.nextToken()
	return true
}

// Errors returns a slice of errors.
func (p *Parser) Errors() []error {
	return p.errors
}

// Parse performs the parsing.
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

	if !p.expectNext(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curr, Value: p.curr.Literal}

	if !p.expectNext(token.ASSIGN) {
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
	if p.trace {
		defer untrace(trace("parseGroupExpression"))
	}
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectNext(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	if p.trace {
		defer untrace(trace("parseIfExpression"))
	}
	exp := &ast.If{Token: p.curr}

	if !p.expectNext(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectNext(token.RPAREN) {
		return nil
	}
	if !p.expectNext(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.next.Type == token.ELSE {
		p.nextToken()

		if !p.expectNext(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.Block {
	if p.trace {
		defer untrace(trace("parseBlockStatement"))
	}
	block := &ast.Block{
		Token:      p.curr,
		Statements: make([]ast.Statement, 0),
	}

	p.nextToken()

	for p.curr.Type != token.RBRACE && p.curr.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	if p.trace {
		defer untrace(trace("parseFunctionLiteral"))
	}
	fn := &ast.Function{Token: p.curr}

	if !p.expectNext(token.LPAREN) {
		return nil
	}
	fn.Args = p.parseFunctionArgs()
	if !p.expectNext(token.LBRACE) {
		return nil
	}
	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseFunctionArgs() []*ast.Identifier {
	if p.trace {
		defer untrace(trace("parseFunctionArgs"))
	}
	args := make([]*ast.Identifier, 0)

	if p.next.Type == token.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, &ast.Identifier{Token: p.curr, Value: p.curr.Literal})

	for p.next.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, &ast.Identifier{Token: p.curr, Value: p.curr.Literal})
	}

	if !p.expectNext(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	if p.trace {
		defer untrace(trace("parseCallExpression"))
	}
	return &ast.Call{
		Token:    p.curr,
		Function: fn,
		Args:     p.parseCallArgs(),
	}
}

func (p *Parser) parseCallArgs() []ast.Expression {
	if p.trace {
		defer untrace(trace("parseCallArgs"))
	}
	args := make([]ast.Expression, 0)

	if p.next.Type == token.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.next.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectNext(token.RPAREN) {
		return nil
	}

	return args
}
