package expr

type nodetype int

const (
	OP nodetype = iota
	VAR
)

//Expr is an algebraic expression
type Expr struct {
	typ nodetype
	sym rune
	l   *Expr
	r   *Expr
}

func (exp Expr) String() string {
	switch exp.typ {
	case OP:
		return "(" + string(exp.sym) + " " + exp.l.String() + " " +
			exp.r.String() + ")"
	case VAR:
		return string(exp.sym)
	default:
		panic("Node of unrecognised type: " + string(exp.typ))
	}
}

//Sub returns a clone of the indexed subexpression. The expressions are indexed
//in preorder
func (exp Expr) Sub(i int) *Expr {
	sub, _ := exp.subrec(i)
	return sub
}

func (exp Expr) subrec(i int) (*Expr, int) {
	if i == 0 {
		return exp.clone(), -1
	}

	if exp.typ == VAR {
		return nil, i
	}
	sub, i := exp.l.subrec(i - 1)
	if sub != nil {
		return sub, -1
	}
	return exp.r.subrec(i - 1)

}

func (exp Expr) clone() *Expr {
	switch exp.typ {
	case VAR:
		return &Expr{exp.typ, exp.sym, nil, nil}
	case OP:
		return &Expr{exp.typ, exp.sym, exp.l.clone(), exp.r.clone()}
	default:
		panic("Node of unrecognised type: " + string(exp.typ))
	}
}
