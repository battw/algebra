package rule

import (
	"github.com/llo-oll/algebra/expr"
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
	_, err := Newrule(lexp, rexp)
	if !wanterror && err != nil {
		t.Errorf("Newrule(%s, %s)\nshould not return an error (%s)", lexp, rexp, msg)
	} else if wanterror && err == nil {
		t.Errorf("Newrule(%s, %s)\nshould return an error (%s)", lexp, rexp, msg)
	}
}
