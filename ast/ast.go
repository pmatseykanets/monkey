package ast

import (
	"bytes"
	"github.com/pmatseykanets/monkey/token"
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
