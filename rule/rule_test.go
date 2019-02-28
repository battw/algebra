package rule

import (
	"github.com/battw/algebra/expr"
	"testing"
)

func Test_Applicable(t *testing.T) {
	applicablise(commute(), "(+ a (* b c))", 2, true, t)
	applicablise(commute(), "(+ a (+ b c))", 2, false, t)
	applicablise(commute(), "(* z E)", 0, true, t)
	applicablise(distrib(), "(& (* (- w e) (+ ($ t r) p)) z)", 1, true, t)
	applicablise(undistrib(), "(£ (+ (* (- c d) q) (* (- c d) r)) s)", 1, true, t)
	applicablise(undistrib(), "(£ (+ (* (- c d) q) (* (- c Z) r)) s)", 1, false, t)
}

func commute() *Rule {
	lhs, _ := expr.Translate("(* a b)")
	rhs, _ := expr.Translate("(* b a)")
	return &Rule{lhs, rhs}
}

func distrib() *Rule {
	lhs, _ := expr.Translate("(* a (+ b c))")
	rhs, _ := expr.Translate("(+ (* a b) (* a c))")
	return &Rule{lhs, rhs}

}

func undistrib() *Rule {
	lhs, _ := expr.Translate("(+ (* a b) (* a c))")
	rhs, _ := expr.Translate("(* a (+ b c))")
	return &Rule{lhs, rhs}
}

func applicablise(r *Rule, expstr string, subi int, desired bool, t *testing.T) {
	exp, err := expr.Translate(expstr)
	if err != nil {
		t.Fatalf("%s\n%s", expstr, err)
	}
	if r.Applicable(exp, subi) != desired {
		if desired {
			t.Errorf("rule\n%s\nshould be applicable to\n%s\n",
				r, exp)
		} else {
			t.Errorf("rule\n%s\nshould NOT be applicable to\n%s\n",
				r, exp)
		}
	}

}

func Test_Apply(t *testing.T) {
	applyalise(commute(), "(* a (* b c))", 0, "(* (* b c) a)", t)
	applyalise(commute(), "(* z E)", 0, "(* E z)", t)
	applyalise(distrib(), "(& (* (- w e) (+ ($ t r) p)) z)", 1,
		"(& (+ (* (- w e) ($ t r)) (* (- w e) p)) z)", t)
	applyalise(undistrib(), "(£ (+ (* (- c d) q) (* (- c d) r)) s)", 1,
		"(£ (* (- c d) (+ q r)) s)", t)

}

func applyalise(r *Rule, expstr string, subi int, desired string, t *testing.T) {
	exp, err := expr.Translate(expstr)
	if err != nil {
		t.Fatalf("%s\n%s", expstr, err)
	}
	desexp, err := expr.Translate(desired)
	if err != nil {
		t.Fatalf("%s\n%s", desired, err)
	}
	res, err := r.Apply(exp, subi)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !res.Equals(desexp) {
		t.Errorf("The output to rule\n%s\non\n%s\nat\n%v\nshould be\n%s\n but is \n%s\n",
			r, exp, subi, desexp, res)
	}
}

func Test_Newrule(t *testing.T) {
	newts("(+ a ($ c d))", "(% (% c d) a)", false, t, "")
	newts("(+ a ($ c d))", "(% a (- c f)", true, t, "extra var on rhs")
	newts("(* z (= x y))", "(+ x x)", false, t, "")

}

func newts(lhs, rhs string, wanterror bool, t *testing.T, msg string) {
	lexp, _ := expr.Translate(lhs)
	rexp, _ := expr.Translate(rhs)
	_, err := New(lexp, rexp)
	if !wanterror && err != nil {
		t.Errorf("rule.New(%s, %s)\nshould not return an error (%s)", lexp, rexp, msg)
	}
}

func Test_Introrule(t *testing.T) {
	exp, err := expr.Translate("(* p (+ (& i j) (* e f)))")
	lhs, err1 := expr.Translate("(+ x y)")
	rhs, err2 := expr.Translate("(* x (* y z))")
	iexp, err3 := expr.Translate("(+ a b)")
	if err != nil || err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("Expressions should translate")
	}
	r, err := New(lhs, rhs)
	if err != nil {
		t.Fatal("Rule should instantiate")
	}
	//A rule including an introduction of a subexp should give an error when
	//no sub exp is provided.
	_, err = r.Apply(exp, 1)
	if err == nil {
		t.Fatalf("Application of rule\n%s\nshould give an error when there "+
			"expression provided to introduce", r)
	}
	//A rule with introduction should give correct results when given good input.
	result, err := r.Apply(exp, 2, iexp)
	if err != nil {
		t.Fatalf("Rule application of\n%s\nshouldn't throw an error\n%s.", r, err)
	}
	desired, err := expr.Translate("(* p (* (& i j) (* (* e f) (+ a b))))")
	if err != nil {
		t.Fatalf("Rule should instantiate")
	}
	if !result.Equals(desired) {
		t.Errorf("result\n%s\nshould be \n%s\n", result, desired)
	}
}

func Test_introductions(t *testing.T) {
	lexp, err1 := expr.Translate("(+ x y)")
	rexp, err2 := expr.Translate("(* x (* y z))")
	if err1 != nil || err2 != nil {
		t.Fatalf("Expressions should translate")
	}
	r, err := New(lexp, rexp)
	if err != nil {
		t.Fatalf("Rule should instantiate")
	}
	intros := r.introductions()
	if len(intros) == 0 {
		t.Errorf("Rule\n%s\nintroduces a variable", r)
	} else if intros[0] != 'z' {
		t.Errorf("Rule\n%s\nintroduces '%c' it should introduce 'z'", r, intros[0])
	}

}

func Test_Pretty(t *testing.T) {
	lhs, err := expr.Translate("(+ a (/ (+ (/ b c) (/ (+ d e) f)) (/ (/ g h) i)))")
	if err != nil {
		t.Fatalf("should translate")
	}
	rhs, err := expr.Translate("(/ (/ (+ a b) c) (/ d (* e f)))")
	if err != nil {
		t.Fatalf("should translate")
	}
	rule1, _ := New(lhs, rhs)
	t.Error(Pretty(rule1))

	rule2, _ := New(rhs, lhs)
	t.Error(Pretty(rule2))

}
