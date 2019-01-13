package expr

import (
	"errors"
	"github.com/llo-oll/algebra/util"
	"strconv"
	"unicode"
)

//TODO Include the brackets as lexical items for better error messages/detection
//Currently, user input can be parsed when it is malformed, meaning that
//an expression is created which might be something other than intended.

//translate, converts strings into expression trees
func Translate(s string) (*Expr, error) {
	runech := util.Runechan(s)
	itemch, err := lex(runech)
	exp, err := parse(itemch)
	return exp, err
}

//item is a lexical symbol
type item struct {
	//TODO define sensible types for these i.e. not strings
	typ nodetype
	sym rune
}

func readop(rch <-chan rune) item {
	sym := <-rch
	if unicode.IsSymbol(sym) {
		return item{OP, sym}
	}
	return item{} //, errors.New(string(sym) + " is not an operator")
}

func readvar(r rune) item {
	return item{VAR, r}
}

func lex(rch <-chan rune) (<-chan item, error) {
	ich := make(chan item)
	go func() {
		for r := range rch {
			switch r {
			case ' ', '\n', '\t', ')':
			case '(':
				item := readop(rch)
				ich <- item
			default:
				ich <- readvar(r)
			}
		}
		close(ich)
	}()
	return ich, nil
}

func parse(ich <-chan item) (*Expr, error) {
	exp, err, size := prec(ich, 0)
	if err == nil {
		item := <-ich
		if item.typ != ERR {
			err = errors.New("Malformed expression: extraneous input '" +
				string(item.sym) + "' at sub " + strconv.Itoa(size+1))
		}
	}
	return exp, err
}
func prec(ich <-chan item, subi int) (*Expr, error, int) {
	i := <-ich
	switch i.typ {
	case OP:
		l, errl, subi := prec(ich, subi+1)
		if errl != nil {
			return l, errl, subi
		}
		r, errr, subi := prec(ich, subi+1)
		if errr != nil {
			return r, errr, subi
		}
		return &Expr{OP, i.sym, l, r}, nil, subi
	case VAR:
		return &Expr{VAR, i.sym, nil, nil}, nil, subi
	default:
		return &Expr{ERR, 949, nil, nil},
			errors.New("Malformed expression: Parse failed at sub expr " +
				strconv.Itoa(subi)),
			subi
	}
}
