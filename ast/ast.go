package ast

import (
	"bytes"
	"github.com/pmatseykanets/monkey/token"
	"strconv"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}

	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var buf bytes.Buffer

	for _, stmt := range p.Statements {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}

// Let represents a let statement.
// E.g. let x = 5;
type Let struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (n *Let) statementNode() {}
func (n *Let) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Let) String() string {
	var buf bytes.Buffer

	buf.WriteString(n.TokenLiteral() + " ")
	buf.WriteString(n.Name.String())
	buf.WriteString(" = ")
	if n.Value != nil {
		buf.WriteString(n.Value.String())
	}
	buf.WriteString(";")
	return buf.String()
}

// Identifier represents an identifier.
// E.g. x in let x = 5;
type Identifier struct {
	Token token.Token
	Value string
}

func (n *Identifier) expressionNode() {}
func (n *Identifier) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Identifier) String() string {
	return n.Value
}

// Return represents return statement.
// E.g. return 5;
// or return foo();
type Return struct {
	Token token.Token // The RETURN token.
	Value Expression
}

func (n *Return) statementNode() {}
func (n *Return) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Return) String() string {
	var buf bytes.Buffer

	buf.WriteString(n.TokenLiteral() + " ")
	if n.Value != nil {
		buf.WriteString(n.Value.String())
	}
	buf.WriteString(";")
	return buf.String()
}

// BareExpr represents a bare Expression statement.
// E.g. x + 10;
type BareExpr struct {
	Token token.Token // The first token of the expression.
	Value Expression
}

func (n *BareExpr) statementNode() {}
func (n *BareExpr) TokenLiteral() string {
	return n.Token.Literal
}
func (n *BareExpr) String() string {
	if n.Value == nil {
		return ""
	}
	return n.Value.String()
}

// IntegerLiteral represents an integer literal.
// E.g. 5;
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (n *IntegerLiteral) expressionNode() {}
func (n *IntegerLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *IntegerLiteral) String() string {
	return strconv.FormatInt(n.Value, 10)
}

// Prefix represents a prefix expression.
// E.g.
// !5
// -10
type Prefix struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (n *Prefix) expressionNode() {}
func (n *Prefix) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Prefix) String() string {
	return "(" + n.Operator + n.Right.String() + ")"
}

// Infix represents an infix expression.
// E.g.
// 5 + 5
// 5 - 5
// 5 * 5
// 5 / 5
// 5 > 5
// 5 < 5
// 5 == 5
// 5 != 5
type Infix struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (n *Infix) expressionNode() {}
func (n *Infix) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Infix) String() string {
	return "(" + n.Left.String() + " " + n.Operator + " " + n.Right.String() + ")"
}

// Boolean represents a boolean literal.
// E.g.
// true
// false
type Boolean struct {
	Token token.Token
	Value bool
}

func (n *Boolean) expressionNode() {}
func (n *Boolean) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Boolean) String() string {
	return strconv.FormatBool(n.Value)
}

// If represents a conditional if expression.
// The else is optional and can be ommited.
// E.g.
// if (x < y) {
// 	return x;
// } else {
// 	return y;
// }
// An if expression produces a value.
// let z = if (x < y) { x } else { y };
type If struct {
	Token       token.Token
	Condition   Expression
	Consequence *Block
	Alternative *Block
}

func (n *If) expressionNode() {}
func (n *If) TokenLiteral() string {
	return n.Token.Literal
}
func (n *If) String() string {
	buf := "if" + n.Condition.String() + " " + n.Consequence.String()
	if n.Alternative == nil {
		return buf
	}

	return buf + "else" + n.Alternative.String()
}

// Block represents a block statement consisting of
// one more statements enslosed in brackets.
// E.g.
// {
// 	a + b;
// }
type Block struct {
	Token      token.Token
	Statements []Statement
}

func (n *Block) statementNode() {}
func (n *Block) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Block) String() string {
	var buf bytes.Buffer

	for _, stmt := range n.Statements {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}

// Function represents a function literal.
// E.g.
// fn(x, y) {
// 	return x + y;
// }
type Function struct {
	Token token.Token
	Args  []*Identifier
	Body  *Block
}

func (n *Function) expressionNode() {}
func (n *Function) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Function) String() string {
	args := make([]string, len(n.Args))
	for i := range n.Args {
		args[i] = n.Args[i].String()
	}

	return n.TokenLiteral() + "(" + strings.Join(args, ", ") + ") " + n.Body.String()
}
