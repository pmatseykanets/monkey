package parser

import (
	"strconv"
	"strings"
	"testing"

	"github.com/pmatseykanets/monkey/ast"
	"github.com/pmatseykanets/monkey/lexer"
)

func TestParseLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	tests := []struct {
		ident string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	lex := lexer.New(strings.NewReader(input))
	p := New(lex)

	prg := p.Parse()
	checkParseErrors(t, p)
	if prg == nil {
		t.Fatal("Program is nil")
	}
	if want, got := 3, len(prg.Statements); want != got {
		t.Fatalf("Expected number of statements %d got %d", want, got)
	}

	for i, tt := range tests {
		stmt := prg.Statements[i]
		if !testLetStatement(t, stmt, tt.ident) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if want, got := "let", stmt.TokenLiteral(); want != got {
		t.Errorf("Expected token literal %s got %s", want, got)
	}

	letStmt, ok := stmt.(*ast.Let)
	if !ok {
		t.Errorf("Expected type *ast.Let got %T", stmt)
		return false
	}

	if want, got := name, letStmt.Name.Value; want != got {
		t.Errorf("Expected letStmt.Name.Value %s got %s", want, got)
		return false
	}

	if want, got := name, letStmt.Name.TokenLiteral(); want != got {
		t.Errorf("Expected letStmt.Name.TokenLiteral() %s got %s", want, got)
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	for _, err := range errors {
		t.Errorf("%s", err)
	}
	t.FailNow()
}

func TestParseReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 5 + 10;
`

	p := New(lexer.FromString(input))

	prg := p.Parse()
	checkParseErrors(t, p)
	if prg == nil {
		t.Fatal("Program is nil")
	}
	if want, got := 3, len(prg.Statements); want != got {
		t.Fatalf("Expected number of statements %d got %d", want, got)
	}

	for _, stmt := range prg.Statements {
		returnStmt, ok := stmt.(*ast.Return)
		if !ok {
			t.Errorf("Expected *ast.Return got %T", stmt)
			continue
		}

		if want, got := "return", returnStmt.TokenLiteral(); want != got {
			t.Errorf("Expected token literal %s got %s", want, got)
		}
	}
}

func TestParseIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := New(lexer.FromString(input))

	prg := p.Parse()
	checkParseErrors(t, p)
	if prg == nil {
		t.Fatal("Program is nil")
	}
	if want, got := 1, len(prg.Statements); want != got {
		t.Fatalf("Expected number of statements %d got %d", want, got)
	}

	stmt, ok := prg.Statements[0].(*ast.BareExpr)
	if !ok {
		t.Fatalf("Expected *ast.BareExpr got %T", prg.Statements[0])
	}

	ident, ok := stmt.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected *ast.Identifier got %T", stmt.Value)
	}
	if want, got := "foobar", ident.Value; want != got {
		t.Errorf("Expected ident.Value %s got %s", want, got)
	}
	if want, got := "foobar", ident.TokenLiteral(); want != got {
		t.Errorf("Expected ident.TokenLiteral %s got %s", want, got)
	}
}

func TestParseIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	p := New(lexer.FromString(input))

	prg := p.Parse()
	checkParseErrors(t, p)
	if prg == nil {
		t.Fatal("Program is nil")
	}
	if want, got := 1, len(prg.Statements); want != got {
		t.Fatalf("Expected number of statements %d got %d", want, got)
	}

	stmt, ok := prg.Statements[0].(*ast.BareExpr)
	if !ok {
		t.Fatalf("Expected *ast.BareExpr got %T", prg.Statements[0])
	}

	literal, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected *ast.IntegerLiteral got %T", stmt.Value)
	}
	if want, got := int64(5), literal.Value; want != got {
		t.Errorf("Expected literal.Value %d got %d", want, got)
	}
	if want, got := "5", literal.TokenLiteral(); want != got {
		t.Errorf("Expected literal.TokenLiteral %s got %s", want, got)
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		p := New(lexer.FromString(tt.input))
		prg := p.Parse()
		checkParseErrors(t, p)

		if prg == nil {
			t.Fatal("Program is nil")
		}
		if want, got := 1, len(prg.Statements); want != got {
			t.Fatalf("Expected number of statements %d got %d", want, got)
		}

		stmt, ok := prg.Statements[0].(*ast.BareExpr)
		if !ok {
			t.Fatalf("Expected *ast.BareExpr got %T", prg.Statements[0])
		}

		pre, ok := stmt.Value.(*ast.Prefix)
		if !ok {
			t.Fatalf("Expected *ast.Prefix got %T", stmt.Value)
		}
		if want, got := tt.operator, pre.Operator; want != got {
			t.Errorf("Expected Operator %s got %s", want, got)
		}
		if !testIntegerLiteral(t, pre.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integer, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected *ast.IntegerLiteral got %T", exp)
		return false
	}

	if want, got := value, integer.Value; want != got {
		t.Errorf("Expected Value %d got %d", want, got)
		return false
	}

	if want, got := strconv.FormatInt(value, 10), integer.TokenLiteral(); want != got {
		t.Errorf("Expected TokenLiteral %s got %s", want, got)
		return false
	}

	return true
}

func TestParseInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {
		p := New(lexer.FromString(tt.input))
		prg := p.Parse()
		checkParseErrors(t, p)

		if prg == nil {
			t.Fatal("Program is nil")
		}
		if want, got := 1, len(prg.Statements); want != got {
			t.Fatalf("Expected number of statements %d got %d", want, got)
		}

		stmt, ok := prg.Statements[0].(*ast.BareExpr)
		if !ok {
			t.Fatalf("Expected *ast.BareExpr got %T", prg.Statements[0])
		}

		exp, ok := stmt.Value.(*ast.Infix)
		if !ok {
			t.Fatalf("Expected *ast.Infix got %T", stmt.Value)
		}
		if !testIntegerLiteral(t, exp.Left, tt.left) {
			return
		}
		if want, got := tt.operator, exp.Operator; want != got {
			t.Errorf("Expected Operator %s got %s", want, got)
		}
		if !testIntegerLiteral(t, exp.Right, tt.right) {
			return
		}
	}
}

func TestParseOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		p := New(lexer.FromString(tt.input))
		prg := p.Parse()
		checkParseErrors(t, p)

		if prg == nil {
			t.Fatal("Program is nil")
		}
		if got := prg.String(); tt.want != got {
			t.Errorf("Expected %s got %s", tt.want, got)
		}
	}
}
