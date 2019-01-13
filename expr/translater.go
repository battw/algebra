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
//OR MAYBE SEPARATE SYNTAX CHECKING FROM PARSING!!!

//translate, converts strings into expression trees
func Translate(s string) (*Expr, error) {
	runech := util.Runechan(s)
	itemch := lex(runech)
	exp, err := parse(itemch)
	return exp, err
}

type itemtype int

const (
	NIL_ITEM itemtype = iota
	ERR_ITEM
	OP_ITEM
	VAR_ITEM
	LBRAK_ITEM
	RBRAK_ITEM
)

//item is a lexical symbol
type item struct {
	typ itemtype
	sym rune
	err error
}

func isOp(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r)
}

func isVar(r rune) bool {
	return unicode.IsLetter(r)
}

func lex(rch <-chan rune) <-chan item {
	ich := make(chan item)
	go func() {
		for r := range rch {
			switch {
			case unicode.IsSpace(r):
			case '(' == r:
			case ')' == r:
			case isOp(r):
				ich <- item{OP_ITEM, r, nil}
			case isVar(r):
				ich <- item{VAR_ITEM, r, nil}
			default:
				ich <- item{
					ERR_ITEM,
					r,
					errors.New("Unrecognised character '" + string(r) + "'")}
			}
		}
		close(ich)
	}()
	return ich
}

func parse(ich <-chan item) (*Expr, error) {
	//TODO size is wrong
	exp, err, size := parserec(ich, 0)
	if err == nil {
		item := <-ich
		if item.typ != NIL_ITEM {
			err = errors.New("Malformed expression: extraneous input '" +
				string(item.sym) + "' at sub " + strconv.Itoa(size+1))
		}
	}
	return exp, err
}
func parserec(ich <-chan item, subi int) (*Expr, error, int) {
	i := <-ich
	switch i.typ {
	case OP_ITEM:
		l, errl, subi := parserec(ich, subi+1)
		if errl != nil {
			return l, errl, subi
		}
		r, errr, subi := parserec(ich, subi+1)
		if errr != nil {
			return r, errr, subi
		}
		return &Expr{OP, i.sym, l, r}, nil, subi
	case VAR_ITEM:
		return &Expr{VAR, i.sym, nil, nil}, nil, subi
	case ERR_ITEM:
		return nil,
			i.err,
			subi
	default:
		return nil,
			errors.New("Malformed expression: something missing"),
			subi
	}
}
