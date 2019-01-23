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
	return r.lhs.Match(e.Subexp(subi))
}

func (r *Rule) Apply(e *expr.Expr, subi int) (*expr.Expr, error) {
	if !r.Applicable(e, subi) {
		return nil, errors.New("Rule\n%s\nis not applicable to \n%s\n",
			r, expr.Subexp(subi))
	}

	return nil, nil
}

func (r *Rule) String() string {
	return r.lhs.String() + " >> " + r.rhs.String()
}
