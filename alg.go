package main

import (
	"bufio"
	"fmt"
	"os"
)

const OPTYPE = "OP"
const VARTYPE = "VAR"

//expr is an algebraic expression
type expr struct {
	typ string //VAR or OP
	sym rune
	l   *expr
	r   *expr
}

func (exp expr) String() string {
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

func parse(sch <-chan item) *expr {
	for s := range sch {
		switch s.typ {
		case OPTYPE:
			return &expr{OPTYPE, s.sym, parse(sch), parse(sch)}
		case VARTYPE:
			return &expr{VARTYPE, s.sym, nil, nil}
		}
	}
	return &expr{}
}

//translate, converts strings into expression trees
func translate(s string) *expr {
	runech := strstream(s)
	itemch := lex(runech)
	return parse(itemch)
}

//eng is the algebra engine.
func eng() (chan<- string, <-chan string) {
	engin := make(chan string)
	engout := make(chan string)
	var exp *expr
	go func() {
	    for instr := range engin {
			exp = translate(instr)
			engout <- exp.String()
		}
	}()
	return engin, engout
}

func main() {
	engin, engout := eng()
	repl(engin, engout)
}
