package expr

import (
	"testing"
)

func TestSimpleParse(t *testing.T) {
	exprs := []string{
		"(+ a b)",
		"(+ (+ a b) c)",
		"(* (- (+ a b) ($ a f)) (# f A))",
	}

	for i, s := range exprs {
		exp, _ := Translate(s)
		if exp.String() != exprs[i] {
			t.Fatalf("Output %s does not match input %s", exp.String(),
				exprs[i])
		}
	}
}

func TestParseFail(t *testing.T) {
	exprs := []string{
		"(a b c)",
		"(+ a b)($ c d)",
		"(+ a)",
		"(; r",
	}

	for i, s := range exprs {
		_, err := Translate(s)
		if err == nil {
			t.Fatalf("Input %s should throw an error", exprs[i])
		}
	}
}
