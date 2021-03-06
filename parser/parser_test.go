package parser

import (
	"strconv"
	"testing"

	"github.com/pmatseykanets/monkey/ast"
	"github.com/pmatseykanets/monkey/lexer"
)

func TestParseLetStatement(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {

		lex := lexer.FromString(tt.input)
		p := New(lex)

		prg := p.Parse()
		checkParseErrors(t, p)
		if prg == nil {
			t.Fatal("Program is nil")
		}
		if want, got := 1, len(prg.Statements); want != got {
			t.Fatalf("Expected statements %d got %d", want, got)
		}

		stmt := prg.Statements[0]
		if !testLetStatement(t, stmt, tt.ident) {
			return
		}
		val := stmt.(*ast.Let).Value
		if !testLiteralExpression(t, val, tt.value) {
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
	tests := []struct {
		input string
		value interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		p := New(lexer.FromString(tt.input))

		prg := p.Parse()
		checkParseErrors(t, p)
		if prg == nil {
			t.Fatal("Program is nil")
		}
		if want, got := 1, len(prg.Statements); want != got {
			t.Fatalf("Expected statements %d got %d", want, got)
		}

		stmt, ok := prg.Statements[0].(*ast.Return)
		if !ok {
			t.Errorf("Expected *ast.Return got %T", prg.Statements[0])
			continue
		}
		if want, got := "return", stmt.TokenLiteral(); want != got {
			t.Errorf("Expected token literal %s got %s", want, got)
		}
		if !testLiteralExpression(t, stmt.Value, tt.value) {
			return
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

	if !testIdentifier(t, stmt.Value, "foobar") {
		return
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected *ast.Identifier got %T", exp)
		return false
	}
	if want, got := value, ident.Value; want != got {
		t.Errorf("Expected ident.Value %s got %s", want, got)
		return false
	}
	if want, got := value, ident.TokenLiteral(); want != got {
		t.Errorf("Expected ident.TokenLiteral %s got %s", want, got)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, want interface{}) bool {
	switch v := want.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("Unhandled exp type %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infix, ok := exp.(*ast.Infix)
	if !ok {
		t.Fatalf("Expected *ast.Infix got %T", exp)
	}
	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}
	if want, got := operator, infix.Operator; want != got {
		t.Errorf("Expected Operator %s got %s", want, got)
		return false
	}
	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
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

	if !testLiteralExpression(t, stmt.Value, 5) {
		return
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
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
		if !testLiteralExpression(t, pre.Right, tt.value) {
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
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !testInfixExpression(t, stmt.Value, tt.left, tt.operator, tt.right) {
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
		{"true", "true"},
		{"false", "false"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
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

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("Expected *ast.Boolean got %T", exp)
		return false
	}

	if want, got := value, boolean.Value; want != got {
		t.Errorf("Expected Value %v got %v", want, got)
		return false
	}

	if want, got := strconv.FormatBool(value), boolean.TokenLiteral(); want != got {
		t.Errorf("Expected TokenLiteral %s got %s", want, got)
		return false
	}

	return true
}

func TestParseBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
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

		if !testLiteralExpression(t, stmt.Value, tt.want) {
			return
		}
	}
}

func TestParseIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

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

	exp, ok := stmt.Value.(*ast.If)
	if !ok {
		t.Fatalf("Expected *ast.If got %T", stmt.Value)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if want, got := 1, len(exp.Consequence.Statements); want != got {
		t.Errorf("Expected statements %d got %d", want, got)
	}

	con, ok := exp.Consequence.Statements[0].(*ast.BareExpr)
	if !ok {
		t.Fatalf("Expected ast.BareExpr statement got %t", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, con.Value, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("Expected Alternative nil got %+v", exp.Alternative)
	}
}

func TestParseIfElseExpression(t *testing.T) {
	input := `
	if (x < y) { 
		x 
	} else {
		y
	}
	`

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

	exp, ok := stmt.Value.(*ast.If)
	if !ok {
		t.Fatalf("Expected *ast.If got %T", stmt.Value)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if want, got := 1, len(exp.Consequence.Statements); want != got {
		t.Errorf("Expected Consequence statements %d got %d", want, got)
	}
	con, ok := exp.Consequence.Statements[0].(*ast.BareExpr)
	if !ok {
		t.Fatalf("Expected ast.BareExpr statement got %t", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, con.Value, "x") {
		return
	}

	if want, got := 1, len(exp.Alternative.Statements); want != got {
		t.Errorf("Expected Alternative statements %d got %d", want, got)
	}
	alt, ok := exp.Alternative.Statements[0].(*ast.BareExpr)
	if !ok {
		t.Fatalf("Expected ast.BareExpr statement got %t", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alt.Value, "y") {
		return
	}
}

func TestParseFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

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

	fn, ok := stmt.Value.(*ast.Function)
	if !ok {
		t.Fatalf("Expected *ast.Function got %T", stmt.Value)
	}

	if want, got := 2, len(fn.Args); want != got {
		t.Fatalf("Expected args %d got %d", want, got)
	}

	testLiteralExpression(t, fn.Args[0], "x")
	testLiteralExpression(t, fn.Args[1], "y")

	if want, got := 1, len(fn.Body.Statements); want != got {
		t.Fatalf("Expected fn body statements %d got %d", want, got)
	}

	bd, ok := fn.Body.Statements[0].(*ast.BareExpr)
	if !ok {
		t.Fatalf("Expected fn body statement type *ast.BareExp got %T", fn.Body.Statements[0])
	}

	testInfixExpression(t, bd.Value, "x", "+", "y")
}

func TestParseFunctionArgs(t *testing.T) {
	tests := []struct {
		input string
		args  []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		p := New(lexer.FromString(tt.input))

		prg := p.Parse()
		checkParseErrors(t, p)
		if prg == nil {
			t.Fatal("Program is nil")
		}

		stmt := prg.Statements[0].(*ast.BareExpr)
		fn := stmt.Value.(*ast.Function)

		if want, got := len(tt.args), len(fn.Args); want != got {
			t.Fatalf("Expected args %d got %d", want, got)
		}

		for i, ident := range tt.args {
			testLiteralExpression(t, fn.Args[i], ident)
		}
	}
}

func TestParseCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

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

	call, ok := stmt.Value.(*ast.Call)
	if !ok {
		t.Fatalf("Expected *ast.Call got %T", stmt.Value)
	}

	if !testIdentifier(t, call.Function, "add") {
		return
	}
	if want, got := 3, len(call.Args); want != got {
		t.Fatalf("Expected args %d got %d", want, got)
	}

	testLiteralExpression(t, call.Args[0], 1)
	testInfixExpression(t, call.Args[1], 2, "*", 3)
	testInfixExpression(t, call.Args[2], 4, "+", 5)
}

func TestParseCallExpressionArgs(t *testing.T) {
	tests := []struct {
		input string
		ident string
		args  []string
	}{
		{"add();", "add", []string{}},
		{"add(1);", "add", []string{"1"}},
		{"add(1, 2 * 3, 4 + 5);", "add", []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tt := range tests {
		p := New(lexer.FromString(tt.input))
		prg := p.Parse()
		checkParseErrors(t, p)
		if prg == nil {
			t.Fatal("Program is nil")
		}

		stmt := prg.Statements[0].(*ast.BareExpr)
		call, ok := stmt.Value.(*ast.Call)
		if !ok {
			t.Fatalf("Expected *ast.Call got %T", stmt.Value)
		}

		if !testIdentifier(t, call.Function, tt.ident) {
			return
		}

		if want, got := len(tt.args), len(call.Args); want != got {
			t.Fatalf("Expected agrs %d got=%d", want, got)
		}

		for i, want := range tt.args {
			if got := call.Args[i].String(); want != got {
				t.Errorf("Expected %d argument %s got %s", i, want, got)
			}
		}
	}
}
