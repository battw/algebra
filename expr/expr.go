package expr

const OPTYPE = "OP"
const VARTYPE = "VAR"

//expr is an algebraic expression
type Expr struct {
	typ string //VAR or OP
	sym rune
	l   *Expr
	r   *Expr
}

func (exp Expr) String() string {
	switch exp.typ {
	case OPTYPE:
		return "(" + string(exp.sym) + " " + exp.l.String() + " " +
			exp.r.String() + ")"
	case VARTYPE:
		return string(exp.sym)
	default:
		return "Bad node"
	}
}

