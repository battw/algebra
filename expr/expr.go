package expr

const OPTYPE = "OP"
const VARTYPE = "VAR"

//Expr is an algebraic expression
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
        panic("Node of unrecognised type: " + exp.typ)
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

    if exp.typ == VARTYPE {
        return nil, i
    }
    sub, i := exp.l.subrec(i-1)
    if sub != nil {
        return sub, -1
    }
    return exp.r.subrec(i-1)

}

func (exp Expr) clone() *Expr {
    switch exp.typ {
    case VARTYPE:
        return &Expr{exp.typ, exp.sym, nil, nil}
    case OPTYPE:
        return &Expr{exp.typ, exp.sym, exp.l.clone(), exp.r.clone()}
    default:
        panic("Node of unrecognised type: " + exp.typ)
    }
}

