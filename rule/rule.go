package rule

import (
	"errors"
	"github.com/llo-oll/algebra/expr"
)

type Rule struct {
	lhs *expr.Expr
	rhs *expr.Expr
}

func Newrule(lhs, rhs *expr.Expr) *Rule {
	return &Rule{lhs, rhs}
}

//Applicable - is the rule applicable to the designated sub expression of the give expression.
func (r *Rule) Applicable(e *expr.Expr, subi int) bool {
	mch, _ := r.lhs.Match(e.Subexp(subi))
	return mch
}

func (r *Rule) Apply(exp *expr.Expr, subi int) (*expr.Expr, error) {
	//Associate variables in the rules lhs with subtrees of the sub expression
	match, varmap := r.lhs.Match(exp.Subexp(subi))
	if !match {
		return nil, errors.New("Rule\n" + r.String() +
			"\nis not applicable to \n" + exp.Subexp(subi).String())
	}
	//Make a copy of the expression
	exp = exp.Clone()
	//Find the subexpression to which the rule is being applied

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
