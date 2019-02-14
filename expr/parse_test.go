package expr

import (
	"testing"
)

func TestSimpleParse(t *testing.T) {
	exprs := []string{
		"(+ a b)",
		"(+ (+ a b) c)",
		"(* (- (+ a b) ($ a f)) (# f A))",
		"(* 1 x)",
		"(& r (- 1 3))",
	}

	for i, s := range exprs {
		exp, err := Translate(s)
		if err != nil {
			t.Fatalf("%s for input %s", err, s)
		}
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
		"(1 a b)",
	}

	for i, s := range exprs {
		_, err := Translate(s)
		if err == nil {
			t.Fatalf("Input %s should return an error", exprs[i])
		}
	}
}
