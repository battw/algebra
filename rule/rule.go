package rule

import (
	"errors"
	"github.com/battw/algebra/expr"
	"strconv"
)

type Rule struct {
	lhs *expr.Expr
	rhs *expr.Expr
}

func New(lhs, rhs *expr.Expr) (*Rule, error) {
	//Test for extra vars on rhs
	return &Rule{lhs, rhs}, nil
}

//Applicable - is the rule applicable to the designated sub expression of the given expression.
func (r *Rule) Applicable(e *expr.Expr, subi int) bool {
	mch, _ := r.lhs.Match(e.Subexp(subi))
	return mch
}

func (r *Rule) Apply(exp *expr.Expr, subi int, introexps ...*expr.Expr) (*expr.Expr, error) {
	//Associate variables in the rules lhs with subtrees of the sub expression
	match, varmap := r.lhs.Match(exp.Subexp(subi))
	if !match {
		return nil, errors.New("Rule\n" + r.String() +
			"\nis not applicable to \n" + exp.Subexp(subi).String())
	}
	var intros []rune = r.introductions()
	if len(intros) != len(introexps) {
		return nil, errors.New("Requires " + strconv.Itoa(len(intros)) + " new expressions." +
			"Got " + strconv.Itoa(len(introexps)))
	}
	for i, sym := range intros {
		varmap[sym] = introexps[i]
	}
	//Make a copy of the expression
	exp = exp.Clone()
	//Make a copy of rhs
	rhs := r.rhs.Clone()
	//Replace the variables in the copy of rhs with the associated subtrees
	rhs = rhs.Subvar(varmap)
	//Substitute the modified rhs back into the copied expression
	exp = exp.Substitute(subi, rhs)
	return exp, nil
}

func (r *Rule) String() string {
	return r.lhs.String() + " -> " + r.rhs.String()
}

//introsvar - If the rule introduces new subexpressions, return the variable symbol
//representing those subexpressions.
//Otherwise return an empty slice.
func (r *Rule) introductions() []rune {
	syms := make([]rune, 0)
	var lvars, rvars map[rune]bool = r.lhs.Vars(), r.rhs.Vars()
	for r, b := range rvars {
		if b && !lvars[r] {
			syms = append(syms, r)
		}
	}
	return syms
}
