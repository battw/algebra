package expr

import (
	"testing"
)

func TestSubstitute(t *testing.T) {
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

func TestSubexp(t *testing.T) {
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
