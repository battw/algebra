package expr

import (
	"testing"
)

func Test_Substitute(t *testing.T) {
	str := "(+ a (* c d))"
	substr := "(+ (%x y) z)"
	exp, _ := Translate(str)
	subexp, _ := Translate(substr)
	result := exp.Substitute(3, subexp)
	if str != exp.String() {
		t.Fatal("exp.Substitute mutates its input")
	}
	desired := "(+ a (* (+ (% x y) z) d))"
	if result.String() != desired {
		t.Fatalf("exp.Substitute failed.\n Got %s\n should be %s\n",
			result.String(), desired)
	}
}

func Test_Subexp(t *testing.T) {
	str := "(+ ($ (% T Q) R) W)"
	subi := 2
	exp, _ := Translate(str)
	sub := exp.Subexp(subi)
	//Test that
	//  it is non destructive
	if exp.String() != str {
		t.Fatalf("exp.Subexp() mutates its input")
	}
	//  it returns the correct result
	desired := "(% T Q)"
	if sub.String() != desired {
		t.Fatalf("exp.Subexp() failed.\n Got %s\n should be %s\n",
			sub.String(), desired)
	}
}

func Test_Equals(t *testing.T) {
	exp1, _ := Translate("(& (+ (& a b) (* r (^ x Y))) (+ o (- r G)))")
	exp2, _ := Translate("(& (+ (& a b) (* r (^ x Y))) (+ o (- r G)))")
	//Basic equals
	if !exp1.Equals(exp2) {
		t.Errorf("\n%s\n%s\nshould evaluate as Equal", exp1, exp2)
	}
	exp3, _ := Translate("(& (+ (& a b) (* r (^ x Y))) (+ t (- r G)))")
	//Differ by a variable name
	if exp1.Equals(exp3) {
		t.Errorf("\n%s\n%s\nshouldn't evaluate as Equal", exp1, exp3)
	}
	exp4, _ := Translate("(& (+ (& a b) (* c d)) (+ e f))")
	//Differ by structure
	if exp1.Equals(exp4) {
		t.Errorf("\n%s\n%s\nshouldn't evaluate as Equal", exp1, exp3)
	}
	//Differ by structure, different case
	if exp4.Equals(exp3) {
		t.Errorf("\n%s\n%s\nshouldn't evaluate as Equal", exp4, exp3)
	}
	exp5, _ := Translate("(* (/ a (+ b c)) (# Q Z))")
	exp6, _ := Translate("(* (/ a b) (# Q Z))")
	//Subtrees don't equal supertrees
	if exp6.Equals(exp5) {
		t.Errorf("\n%s\n%s\nAren't equal", exp6, exp5)
	}
	//Supertrees don't equal subtrees
	if exp5.Equals(exp6) {
		t.Errorf("\n%s\n%s\nAren't equal", exp5, exp6)
	}
}
func Test_Match(t *testing.T) {
	sub, _ := Translate("(* (- a b) (* c d))")
	sup, _ := Translate("(* (- z (^ c d)) (* (+ e f) (+ g h))")
	if !sub.Match(sup) {
		t.Errorf("\n%s\n%s\nshould 'Match'", sub, sup)
	}
	sub, _ = Translate("(* (- a b) (* c d))")
	sup, _ = Translate("(* (- z (^ c d)) (- (+ e f) (+ g h))")
	if sub.Match(sup) {
		t.Errorf("\n%s\n%s\nshouldn't 'Match'", sub, sup)
	}
	sub, _ = Translate("(+ a a)")
	sup, _ = Translate("(+ (* x y) (* x y))")
	//Repeated vars left
	if !sub.Match(sup) {
		t.Errorf("\n%s\n%s\nshould 'Match'", sub, sup)
	}
	sup, _ = Translate("(+ (* x y) (* q y))")
	if sub.Match(sup) {
		t.Errorf("\n%s\n%s\nshouldn't 'Match'", sub, sup)
	}

}
