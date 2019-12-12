package eval

import (
	"testing"

	"github.com/pmatseykanets/monkey/lexer"
	"github.com/pmatseykanets/monkey/object"
	"github.com/pmatseykanets/monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.want)
	}
}

func testEval(input string) object.Object {
	p := parser.New(lexer.FromString(input))
	prg := p.Parse()

	return Eval(prg)
}

func testIntegerObject(t *testing.T, obj object.Object, want int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Expected object.Integer got %T (%v)", obj, obj)
		return false
	}
	if got := result.Value; want != got {
		t.Errorf("Expected Value %d got %d", want, got)
		return false
	}

	return true
}
