package parser

import (
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

func TestIdentifierExpression(t *testing.T) {
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
		t.Errorf("Expected *ast.Identifier got %T", stmt.Value)
	}
	if want, got := "foobar", ident.Value; want != got {
		t.Errorf("Expected ident.Value %s got %s", want, got)
	}
	if want, got := "foobar", ident.TokenLiteral(); want != got {
		t.Errorf("Expected ident.TokenLiteral %s got %s", want, got)
	}
}
