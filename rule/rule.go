package rule

import (
	"github.com/llo-oll/algebra/expr"
)

type Rule struct {
	lhs expr.Expr
	rhs expr.Expr
}

func Newrule(lhs, rhs *expr.Expr) *rule {
	Rule{lhs, rhs}
}
