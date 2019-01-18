package main

import (
	"github.com/llo-oll/algebra/toknify"
	"github.com/llo-oll/algebra/util"
	"testing"
)

func Test_paramcheck(t *testing.T) {
	num := toknify.INT
	name := toknify.NAME
	exp := toknify.EXPR
	// nul := toknify.NIL
	// err := toknify.ERR

	sig0 := []toknify.Toktyp{name, num, name, exp}
	sig1 := []toknify.Toktyp{name, name, exp}
	sig2 := []toknify.Toktyp{name}
	sigs := [][]toknify.Toktyp{sig0, sig1, sig2}

	//Returns the correct things for a simple, positive case
	tokch := toknify.Tokenise(util.Runechan("barry harry (+ a b)"))
	toks, sigi, err := paramcheck(sigs, tokch)
	if len(toks) != 3 {
		t.Errorf("Returned the wrong number of tokens, should be 3, was %v\n",
			len(toks))
	}
	if sigi != 1 {
		t.Errorf("Returned the wrong index, should be 1, is %v\n", sigi)
	}
	if err != nil {
		t.Errorf("Returned an error when it shouldn't have: %s\n", err)
	}

	//Returns an error when given too few arguments
	tokch = toknify.Tokenise(util.Runechan("barry harry"))
	toks, sigi, err = paramcheck(sigs, tokch)
	if err == nil {
		t.Error("Didn't return error when given incorrect tokens as args: barry, harry")
	}
	tokch = toknify.Tokenise(util.Runechan(""))
	toks, sigi, err = paramcheck(sigs, tokch)
	if err == nil {
		t.Error("Didn't return error when given no tokens as args")
	}
	tokch = toknify.Tokenise(util.Runechan("baz 12 haz (* b c) 10"))
	toks, sigi, err = paramcheck(sigs, tokch)
	if err == nil {
		t.Error("Didn't return error when given too many tokens as args")
	}

}

func Test_exprdef(t *testing.T) {
	instr := "(+ (* b c) (- a d))"
	varstr := "var"
	tokch := toknify.Tokenise(util.Runechan(varstr + " " + instr))
	env := newenviron()
	outstr := exprdef(tokch, env)
	//Stores correct expression in the map
	if env.expmap[varstr].String() != instr {
		t.Errorf("Stored expression %s does not match input expression %s\n",
			instr, env.expmap[varstr])
	}
	//Returns the appropriate string
	if instr != outstr {
		t.Errorf("The return string %s doesn't match the input expression %s",
			outstr, instr)
	}
}

func Test_subexpr(t *testing.T) {
	instr := "(+ (* b c) (- a d))"
	varstr := "var"
	subi := "5"
	substr := "(- a d)"
	tokch0 := toknify.Tokenise(util.Runechan(varstr + " " + instr))
	env := newenviron()
	exprdef(tokch0, env)
	tokch1 := toknify.Tokenise(util.Runechan(varstr + " " + subi))
	outstr := subexpr(tokch1, env)
	if outstr != substr {
		t.Errorf("The return string %s doesn't match the correct substring %s",
			outstr, substr)
	}

}

func Test_substitute(t *testing.T) {
	expvar := "x"
	expstr := "(+ (@ t (% q r)) s)"
	subvar := "y"
	substr := "(& a b)"
	subi := "6"
	reqstr := "(+ (@ t (% q (& a b))) s)"
	env := newenviron()
	tokch0 := toknify.Tokenise(util.Runechan(expvar + " " + expstr))
	exprdef(tokch0, env)
	tokch1 := toknify.Tokenise(util.Runechan(subvar + " " + substr))
	exprdef(tokch1, env)
	tokch2 := toknify.Tokenise(util.Runechan(expvar + " " + subi + " " + subvar))
	outstr := substitute(tokch2, env)
	if outstr != reqstr {
		t.Errorf("The return string %s doesn't match required string %s",
			outstr, reqstr)
	}
}
