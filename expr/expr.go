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

//item is a lexical symbol in the algebra
type item struct {
	//TODO can I define sensible types for these i.e. not strings
	typ string
	sym rune
}

func strstream(s string) <-chan rune {
	rch := make(chan rune)
	go func() {
		for _, c := range s {
			rch <- c
		}
		close(rch)
	}()
	return rch
}

func readop(rch <-chan rune) item {
	switch <-rch {
	case '+':
		return item{OPTYPE, '+'}
	default:
		panic("Operation is invalid")
	}
}

func readvar(r rune) item {
	return item{"VAR", r}
}

func lex(rch <-chan rune) <-chan item {
	sch := make(chan item)
	go func() {
		for r := range rch {
			switch r {
			case ' ', '\n', '\t', ')':
			case '(':
				sch <- readop(rch)
			default:
				sch <- readvar(r)
			}
		}
		close(sch)
	}()
	return sch
}

func parse(sch <-chan item) *Expr {
	for s := range sch {
		switch s.typ {
		case OPTYPE:
			return &Expr{OPTYPE, s.sym, parse(sch), parse(sch)}
		case VARTYPE:
			return &Expr{VARTYPE, s.sym, nil, nil}
		}
	}
	return &Expr{}
}

//translate, converts strings into expression trees
func Translate(s string) *Expr {
	runech := strstream(s)
	itemch := lex(runech)
	return parse(itemch)
}

