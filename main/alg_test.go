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
