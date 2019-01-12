package expr

//translate, converts strings into expression trees
func Translate(s string) *Expr {
	runech := Strstream(s)
	itemch := lex(runech)
	return parse(itemch)
}

//item is a lexical symbol
type item struct {
	//TODO define sensible types for these i.e. not strings
	typ nodetype
	sym rune
}

//TODO refactor this into a utility package as it is useful
func Strstream(s string) <-chan rune {
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
		return item{OP, '+'}
	default:
		//TODO catch this, or something. Need decent error messages.
		panic("Operation is invalid")
	}
}

func readvar(r rune) item {
	return item{VAR, r}
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
		case OP:
			return &Expr{OP, s.sym, parse(sch), parse(sch)}
		case VAR:
			return &Expr{VAR, s.sym, nil, nil}
		}
	}
	return &Expr{}
}
