package eval

import (
	"testing"

	"github.com/rumpl/monkey-lang/lexer"
	"github.com/rumpl/monkey-lang/object"
	"github.com/rumpl/monkey-lang/parser"
)

func TestEvalIntegerExperssion(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("object is not Integer, got %t", obj)
		return false
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value, got %d, want %d", result.Value, expected)
		return false
	}

	return true
}
