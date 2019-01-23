package expr

import (
	"fmt"
)

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
		return "OP "
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

//Equals - Expression equality: same shape, operations and variables.
func (exp *Expr) Equals(other *Expr) bool {
	if exp.typ == other.typ && exp.sym == other.sym {
		if exp.typ == VAR {
			return true
		} else if exp.typ == OP {
			return exp.l.Equals(other.l) && exp.r.Equals(other.r)
		}
	}
	return false
}

//Match - Checks that 'sub' matches a sub tree of 'sup' with the same root.
//Returns true if 'sup' can be constructed by replacing each of the variables in 'sub'
//with some expression. When variables appear more than once in sub, their corresponding
//subtrees in sup are checked for equivalence.
func (sub *Expr) Match(sup *Expr) bool {
	var2exps := make(map[rune][]*Expr)
	//Check that 'sub' is a subtree of 'sup'
	match := sub.matchrec(sup, var2exps)
	fmt.Println(match)
	//Check for equality of subtrees in 'sup' which are represented by the same variable in 'sub'
	for _, exps := range var2exps {
		if len(exps) > 1 {
			for i := 1; i < len(exps); i++ {
				match = match && exps[i-1].Equals(exps[i])
			}
		}
	}
	return match
}

func (sub *Expr) matchrec(sup *Expr, var2exps map[rune][]*Expr) bool {
	if sub.typ == VAR {
		if var2exps[sub.sym] == nil {
			var2exps[sub.sym] = make([]*Expr, 0, 3)
		}
		var2exps[sub.sym] = append(var2exps[sub.sym], sup)
		return true
	}
	if sup.typ == VAR {
		return false
	}
	if sub.typ == OP && sup.typ == OP && sub.sym == sup.sym {
		return sub.l.matchrec(sup.l, var2exps) && sub.r.matchrec(sup.r, var2exps)
	}
	return false
}

//Sub returns a clone of the indexed subexpression. The expressions are indexed
//in preorder starting at 0
func (exp Expr) Subexp(i int) *Expr {
	sub, _ := exp.subrec(i)
	return sub.clone()
}

//TODO This needs to handle errors (index out of bounds)
//WARNING any changes to the returned sub tree will affect the input tree.
func (exp *Expr) subrec(i int) (*Expr, int) {
	if i == 0 {
		return exp, -1
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
func (exp *Expr) clone() *Expr {
	switch exp.typ {
	case VAR:
		return &Expr{exp.typ, exp.sym, nil, nil}
	case OP:
		return &Expr{exp.typ, exp.sym, exp.l.clone(), exp.r.clone()}
	default:
		return &Expr{ERR, 949, nil, nil}
	}
}

//Substitute returns a new expression where the sub expression at index 'subi'
//is replaced with 'subexp'.
func (exp *Expr) Substitute(subi int, substitute *Expr) *Expr {
	exp = exp.clone()
	old, _ := exp.subrec(subi)
	old.typ = substitute.typ
	old.sym = substitute.sym
	old.l = substitute.l
	old.r = substitute.r
	return exp
}
