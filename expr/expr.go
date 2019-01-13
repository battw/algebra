package expr

type nodetype int

const (
	ERR nodetype = iota
	OP
	VAR
)

func (nt nodetype) String() string {
	switch nt {
	case ERR:
		return "ERR"
	case OP:
		return "OP"
	case VAR:
		return "VAR"
	default:
		return "INVALID"
	}
}

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
		return "ERR: " + string(exp.sym)
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

//TODO return an error rather than badly formed tree
func (exp Expr) clone() *Expr {
	switch exp.typ {
	case VAR:
		return &Expr{exp.typ, exp.sym, nil, nil}
	case OP:
		return &Expr{exp.typ, exp.sym, exp.l.clone(), exp.r.clone()}
	default:
		return &Expr{ERR, 949, nil, nil}
	}
}
