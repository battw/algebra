package main

import (
	"bufio"
	"fmt"
	"github.com/llo-oll/algebra/expr"
	"github.com/llo-oll/algebra/toknify"
	"os"
	"strconv"
)

type environ struct {
	expmap map[string]*expr.Expr
}

func newenviron() *environ {
	expmap := make(map[string]*expr.Expr)
	env := &environ{expmap}
	return env
}

func main() {
	engin, engout := eng()
	repl(engin, engout)
}

func repl(engin chan<- string, engout <-chan string) {
	bio := bufio.NewReader(os.Stdin)
	for {
		//read
		line, _ := bio.ReadString('\n')
		//execute
		engin <- line
		//print
		str := <-engout
		if len(str) > 0 {
			fmt.Println(str)
		}
	}
}

//eng is the algebra engine.
func eng() (chan<- string, <-chan string) {
	engin := make(chan string)
	engout := make(chan string)
	env := newenviron()
	go func() {
		for instr := range engin {
			engout <- handle(instr, env)
		}
	}()
	return engin, engout
}

func handle(input string, env *environ) string {
	runech := expr.Strstream(input)
	tokch := toknify.Tokenise(runech)
	cmdtok := <-tokch
	switch cmdtok.Str {
	case "e", "expr":
		return exprdef(tokch, env)
	case "p", "print":
		return printexp(tokch, env)
	case "pv", "printvars":
		return printvars(tokch, env)
	case "r", "rule":
		return ruledef(tokch, env)
	case "a", "apply":
		return applyrule(tokch, env)
	case "s", "sub":
		return subexpr(tokch, env)
	case "sbs":
		return substitute(tokch, env)
	case "":
		return ""
	default:
		return cmdtok.Str + " is not a command"
	}
}

func typeerr(tok toknify.Tokn, should toknify.Toktyp) string {
	return tok.Str + " is of type " + tok.Typ.String() +
		", should be of type " + should.String()
}

func exprdef(tokch <-chan toknify.Tokn, env *environ) string {
	tok1 := <-tokch
	tok2 := <-tokch
	if tok1.Typ != toknify.NAME {
		return typeerr(tok1, toknify.NAME)
	}
	if tok2.Typ != toknify.EXPR {
		return typeerr(tok2, toknify.EXPR)
	}
	//parse the expression
	exp := expr.Translate(tok2.Str)
	//add to the environment
	env.expmap[tok1.Str] = exp
	return env.expmap[tok1.Str].String()
}

func printexp(tokch <-chan toknify.Tokn, env *environ) string {
	tok := <-tokch
	if tok.Typ != toknify.NAME {
		return typeerr(tok, toknify.NAME)
	}
	exp := env.expmap[tok.Str]
	if exp == nil {
		return "There is no expression named " + tok.Str
	}
	return exp.String()
}

func printvars(tokch <-chan toknify.Tokn, env *environ) string {
	str := ""
	for k, exp := range env.expmap {
		str += k + ": " + exp.String() + "\n"
	}
	return str[:len(str)-1]
}

func ruledef(tokch <-chan toknify.Tokn, env *environ) string {
	return "RULE"
}

func applyrule(tokch <-chan toknify.Tokn, env *environ) string {
	return "APPLY"
}

func subexpr(tokch <-chan toknify.Tokn, env *environ) string {
	tok1 := <-tokch
	if tok1.Typ != toknify.NAME {
		return typeerr(tok1, toknify.NAME)
	}
	tok2 := <-tokch
	if tok2.Typ != toknify.INT {
		return typeerr(tok2, toknify.INT)
	}
	exp := env.expmap[tok1.Str]
	if exp == nil {
		return "There is no expression named " + tok1.Str
	}
	subi, _ := strconv.Atoi(tok2.Str)
	return exp.Sub(subi).String()
}

func substitute(tokch <-chan toknify.Tokn, env *environ) string {
	return "SUBSTITUTE"
}
