package main

import (
	"bufio"
	"fmt"
	"github.com/llo-oll/algebra/expr"
	"os"
	"strconv"
	//	"strings"
	"unicode"
)

type toktyp int

const (
	NIL toktyp = iota
	ERR
	EXPR
	INT
	NAME
)

func (tt toktyp) String() string {
	switch tt {
	case NIL:
		return "NIL"
	case ERR:
		return "ERR"
	case EXPR:
		return "EXPR"
	case INT:
		return "INT"
	case NAME:
		return "NAME"
	default:
		return "INVALID TYPE"
	}
}

type tokn struct {
	typ toktyp
	str string
}

func (tok tokn) String() string {
	return fmt.Sprintf("tokn{%s, %s}", tok.typ, tok.str)
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isName(s string) bool {
	isname := unicode.IsLetter(rune(s[0]))
	for _, r := range s {
		if unicode.IsSpace(r) {
			isname = false
			break
		}
	}
	return isname
}

func isExpr(s string) bool {
	return len(s) > 0 && s[0] == '('
}

func repl(engin chan<- string, engout <-chan string) {
	bio := bufio.NewReader(os.Stdin)
	for {
		//read
		line, _ := bio.ReadString('\n')
		//execute
		engin <- line
		//print
		fmt.Println(<-engout)
	}
}

//eng is the algebra engine.
func eng() (chan<- string, <-chan string) {
	engin := make(chan string)
	engout := make(chan string)
	expmap := make(map[string]*expr.Expr)
	go func() {
		for instr := range engin {
			engout <- handle(instr, expmap)
		}
	}()
	return engin, engout
}

//readtok returns the next token.
func readtokstr(rch <-chan rune) string {
	//ignore any initial whitespace
	var r rune
	for r = range rch {
		if !unicode.IsSpace(r) {
			break
		}
	}

	var tokstr string
	//if starts with a ( then extract the expression string
	if r == '(' {
		tokstr += string(r)
		scope := 1
		for r := range rch {
			switch r {
			case '(':
				scope++
			case ')':
				scope--
			}
			tokstr += string(r)
			if scope == 0 {
				break
			}
		}
		//These are non-token strings
	} else if r == '\000' || unicode.IsSpace(r) {
		tokstr = ""
	} else {
		tokstr += string(r)
		for r = range rch {
			if unicode.IsSpace(r) {
				break
			}
			tokstr += string(r)
		}
	}

	return tokstr
}

func str2tokn(str string) tokn {
	var tok tokn
	tok.str = str
	if len(str) == 0 {
		tok.typ = NIL
	} else if isInt(str) {
		tok.typ = INT
	} else if isName(str) {
		tok.typ = NAME
	} else if isExpr(str) {
		tok.typ = EXPR
	} else {
		tok.typ = ERR
		tok.str = "Token string '" + str + "' isn't a NAME, INT or EXPR"
	}
	return tok
}

//tokenise, returns a channel providing tokens as strings.
//A token is either a space separated command/argument or an expr.
func tokenise(rch <-chan rune) <-chan tokn {
	tokch := make(chan tokn)
	go func() {
		for {
			tokstr := readtokstr(rch)
			tok := str2tokn(tokstr)
			if tok.typ == NIL {
				break
			}
			tokch <- tok
		}
		close(tokch)
	}()
	return tokch
}

func handle(input string, expmap map[string]*expr.Expr) string {
	runech := expr.Strstream(input)
	tokch := tokenise(runech)
	for tok := range tokch {
		fmt.Println(tok)
	}
	return " "
}

func main() {
	engin, engout := eng()
	repl(engin, engout)
}
