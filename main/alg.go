package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/llo-oll/algebra/expr"
	"github.com/llo-oll/algebra/toknify"
	"github.com/llo-oll/algebra/util"
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
	runech := util.Runechan(input)
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
	errmsg := ""
	if tok.Typ == toknify.NIL {
		errmsg += "Missing argument should be of type " + should.String()
	} else {
		errmsg += tok.Str + " is of type " + tok.Typ.String() +
			", should be of type " + should.String()
	}
	return errmsg
}

//paramcheck takes lists of correct token types, and checks them against those coming out of the
//channel. If they are correct, it returns a list containing them, an int giving the index of the
//matched param list.
func paramcheck(
	desired [][]toknify.Toktyp, tokch <-chan toknify.Tokn) ([]toknify.Tokn, int, error) {

	toks := make([]toknify.Tokn, 0, 10)
	//put tokens in slice
	for tok := range tokch {
		toks = append(toks, tok)
	}
	//check slice has a matching number of toks with some desired signature
	contender := make([]bool, len(desired))
	for i := 0; i < len(desired); i++ {
		contender[i] = len(desired[i]) == len(toks)
	}
	//check that something with same number of toks matches
	for j := 0; j < len(toks); j++ {
		for i := 0; i < len(desired); i++ {
			if contender[i] && desired[i][j] != toks[j].Typ {
				contender[i] = false
			}
		}
	}
	// See if there's a match
	matchi := -1
	for i := 0; i < len(desired); i++ {
		if contender[i] {
			matchi = i
			break
		}
	}
	var err error
	if matchi == -1 {
		prmstr := ""
		for _, prm := range toks {
			prmstr += prm.Typ.String() + " "
		}
		prmstr += "\n"
		desstr := ""
		for _, des := range desired {
			for _, typ := range des {
				desstr += typ.String() + " "
			}
			desstr += "\n"
		}
		err = errors.New(fmt.Sprintf("The parameters of type\n%s"+
			"dont't match any of\n%s", prmstr, desstr))
	}
	return toks, matchi, err
}

func exprdef(tokch <-chan toknify.Tokn, env *environ) string {
	//parse the expression
	desired := [][]toknify.Toktyp{{toknify.NAME, toknify.EXPR}}
	toks, _, err := paramcheck(desired, tokch)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	exp, err := expr.Translate(toks[1].Str)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	//add to the environment
	env.expmap[toks[0].Str] = exp
	return env.expmap[toks[0].Str].String()
}

func printexp(tokch <-chan toknify.Tokn, env *environ) string {
	desired := [][]toknify.Toktyp{{toknify.NAME}}
	toks, _, err := paramcheck(desired, tokch)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	exp := env.expmap[toks[0].Str]
	if exp == nil {
		return "There is no expression named " + toks[0].Str
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
	desired := [][]toknify.Toktyp{{toknify.NAME, toknify.INT}}
	toks, _, err := paramcheck(desired, tokch)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	exp := env.expmap[toks[0].Str]
	if exp == nil {
		return "There is no expression named " + toks[0].Str
	}
	subi, _ := strconv.Atoi(toks[1].Str)
	return exp.Subexp(subi - 1).String() // -1 so counting from 1 up rather than 0
}

func substitute(tokch <-chan toknify.Tokn, env *environ) string {
	desired := [][]toknify.Toktyp{{toknify.NAME, toknify.INT, toknify.NAME}}
	toks, _, err := paramcheck(desired, tokch)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	exp := env.expmap[toks[0].Str]
	if exp == nil {
		return "There is no expression named " + toks[0].Str
	}
	subi, _ := strconv.Atoi(toks[1].Str)
	subexp := env.expmap[toks[2].Str]
	if subexp == nil {
		return "There is no expression named " + toks[2].Str
	}
	return exp.Substitute(subi, subexp).String()
}
