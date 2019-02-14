package expr

type nodetype int

const (
	ERR nodetype = iota
	OP
	VAR
	NUM
)

func (nt nodetype) String() string {
	switch nt {
	case ERR:
		return "ERR"
	case OP:
		return "OP "
	case VAR:
		return "VAR"
	case NUM:
		return "NUM"
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
	case VAR, NUM:
		return string(exp.sym)
	default:
		return "ERR: " + string(exp.sym)
	}
}

//Equals - Expression equality: same shape, operations and variables.
func (exp *Expr) Equals(other *Expr) bool {
	if exp.typ == other.typ && exp.sym == other.sym {
		if exp.typ == VAR || exp.typ == NUM {
			return true
		} else if exp.typ == OP {
			return exp.l.Equals(other.l) && exp.r.Equals(other.r)
		}
	}
	return false
}

//Match - Checks that 'sub' matches a sub tree of 'sup' with the same root.
//Returns true if 'sup' can be constructed by replacing each of the variables in
//'sub' with some expression. When variables appear more than once in sub, their
//corresponding subtrees in sup are checked for equivalence.
//Also returns a map from the variables in sub to their corresponding
//subexpressions in sup
func (sub *Expr) Match(sup *Expr) (bool, map[rune]*Expr) {
	var2exps := make(map[rune][]*Expr)
	//Check that 'sub' is a subtree of 'sup'
	match := sub.matchrec(sup, var2exps)
	//Check for equality of subtrees in 'sup' which are represented by the same variable in 'sub'
	for _, exps := range var2exps {
		if len(exps) > 1 {
			for i := 1; i < len(exps); i++ {
				match = match && exps[i-1].Equals(exps[i])
			}
		}
	}
	v2e := make(map[rune]*Expr)
	for k, exps := range var2exps {
		if len(exps) > 0 {
			v2e[k] = exps[0]
		}
	}
	return match, v2e
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
	if sub.typ == NUM && sup.typ == NUM && sub.sym == sup.sym {
		return true
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
	return sub.Clone()
}

//TODO This needs to handle errors (index out of bounds)
//WARNING any changes to the returned sub tree will affect the input tree.
func (exp *Expr) subrec(i int) (*Expr, int) {
	if i == 0 {
		return exp, -1
	}

	if exp.typ == VAR || exp.typ == NUM {
		return nil, i
	}
	sub, i := exp.l.subrec(i - 1)
	if sub != nil {
		return sub, -1
	}
	return exp.r.subrec(i - 1)

}

//TODO return an error rather than badly formed tree
func (exp *Expr) Clone() *Expr {
	switch exp.typ {
	case VAR, NUM:
		return &Expr{exp.typ, exp.sym, nil, nil}
	case OP:
		return &Expr{exp.typ, exp.sym, exp.l.Clone(), exp.r.Clone()}

	default:
		return &Expr{ERR, 949, nil, nil}
	}
}

//Substitute returns a new expression where the sub expression at index 'subi'
//is replaced with 'subexp'.
func (exp *Expr) Substitute(subi int, substitute *Expr) *Expr {
	exp = exp.Clone()
	old, _ := exp.subrec(subi)
	old.typ = substitute.typ
	old.sym = substitute.sym
	old.l = substitute.l
	old.r = substitute.r
	return exp
}

//Subvar - for any occurance of a VAR whose symbol is a key in the given map,
//substitute the Var with the mapped Expr.
//Return the result.
//This operation is non-destructive.
func (exp *Expr) Subvar(varmap map[rune]*Expr) *Expr {
	exp = exp.Clone()
	var rec func(*Expr) *Expr
	rec = func(sub *Expr) *Expr {
		if sub.typ == VAR && varmap[sub.sym] != nil {
			return varmap[sub.sym].Clone()
		}
		if sub.typ == OP {
			sub.l = rec(sub.l)
			sub.r = rec(sub.r)
		}
		return sub
	}
	exp = rec(exp)

	return exp
}

//Vars - returns the set of variable symbols from the expression
func (exp *Expr) Vars() map[rune]bool {
	symset := make(map[rune]bool)
	var rec func(*Expr)
	rec = func(sub *Expr) {
		if sub.typ == VAR {
			symset[sub.sym] = true
			return
		}
		if sub.typ == OP {
			rec(sub.l)
			rec(sub.r)
		}
	}
	rec(exp)
	return symset
}
