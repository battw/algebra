package main

import (
	"github.com/battw/algebra/toknify"
	"github.com/battw/algebra/util"
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

func helper_defexps(defstrs []string, env *environ) {
	for _, str := range defstrs {
		tokch := toknify.Tokenise(util.Runechan(str))
		exprdef(tokch, env)
	}
}

//commuterule adds a commute rule (* a b) -> (* b a) to the environment.
//desired and actual string outputs plus the rule name.
func helper_commuterule(env *environ) (string, string, string) {
	rname := "commute"
	preexp := "(* a b)"
	postexp := "(* b a)"
	rtokch := toknify.Tokenise(util.Runechan(rname + " " + preexp + " " + postexp))
	outstr := ruledef(rtokch, env)
	return preexp + " -> " + postexp, outstr, rname
}

func helper_distribrule(env *environ) (string, string, string) {
	rname := "distrb"
	preexp := "(* a (+ b c))"
	postexp := "(+ (* a b) (* a c))"
	rtokch := toknify.Tokenise(util.Runechan(rname + " " + preexp + " " + postexp))
	outstr := ruledef(rtokch, env)
	return preexp + " -> " + postexp, outstr, rname
}

func helper_undistribrule(env *environ) (string, string, string) {
	rname := "undistrb"
	preexp := "(+ (* a b) (* a c))"
	postexp := "(* a (+ b c))"
	rtokch := toknify.Tokenise(util.Runechan(rname + " " + preexp + " " + postexp))
	outstr := ruledef(rtokch, env)
	return preexp + " -> " + postexp, outstr, rname

}

func Test_ruledef(t *testing.T) {
	env := newenviron()
	desiredstr, actualstr, _ := helper_commuterule(env)
	if desiredstr != actualstr {
		t.Errorf("The rule string\n%s\n"+
			"doesn't match the desired commute rule string\n%s\n",
			actualstr, desiredstr)
	}

	env = newenviron()
	desiredstr, actualstr, _ = helper_distribrule(env)
	if desiredstr != actualstr {
		t.Errorf("The rule string\n%s\n"+
			"doesn't match the desired distribute rule string\n%s\n",
			actualstr, desiredstr)
	}

	env = newenviron()
	desiredstr, actualstr, _ = helper_undistribrule(env)
	if desiredstr != actualstr {
		t.Errorf("The rule string\n%s\n"+
			"doesn't match the desired undistribute rule string\n%s\n",
			actualstr, desiredstr)
	}

}

func helper_applyrule(rname string, ename string, i string, resname string, env *environ) string {
	tokch := toknify.Tokenise(util.Runechan(rname + " " + ename + " " + i + " " + resname))
	return applyrule(tokch, env)
}

func Test_applyrule(t *testing.T) {
	//Get error from rule which doesn't match the specified sub expression
	env := newenviron()
	_, _, rname := helper_commuterule(env)
	ename := "ex"
	exp := "(* (* (* a b) c) (* e f))"
	i := "2"
	helper_defexps([]string{ename + " " + exp}, env)
	result := helper_applyrule(rname, ename, i, "", env)
	if len(result) == 0 || result[0] != '(' {
		t.Errorf("Rule application should return an error string " +
			"as the rule isn't applicable")
	}

	//Apply a rule at root of expression
	env = newenviron()
	_, _, rname = helper_commuterule(env)
	ename = "ex"
	exp = "(* (* a b) (* c d))"
	i = "1"
	desired := "(* (* c d) (* a b))"
	helper_defexps([]string{ename + " " + exp}, env)
	result = helper_applyrule(rname, ename, i, "", env)
	if desired != result {
		t.Errorf("The result of apply commutative rule at %v is\n%s\n"+
			"It should be \n%s\n", i, result, desired)
	}
	//Apply rule deeper than at root
	env = newenviron()
	_, _, rname = helper_commuterule(env)
	ename = "ex"
	exp = "(* (* (* a b) c) (* e f))"
	i = "2"
	desired = "(* (* c (* a b)) (* e f))"
	helper_defexps([]string{ename + " " + exp}, env)
	result = helper_applyrule(rname, ename, i, "", env)
	if desired != result {
		t.Errorf("The result of apply commutative rule at %v is\n%s\n"+
			"It should be \n%s\n", i, result, desired)
	}
	//Apply rule where one var in the lhs appears twice in the rhs
	env = newenviron()
	_, _, rname = helper_distribrule(env)
	ename = "ex"
	exp = "(* a (+ (+ c d) e))"
	i = "1"
	desired = "(+ (* a (+ c d)) (* a e))"
	helper_defexps([]string{ename + " " + exp}, env)
	result = helper_applyrule(rname, ename, i, "", env)
	if desired != result {
		t.Errorf("The result of apply distributive rule at %v is\n%s\n"+
			"It should be \n%s\n", i, result, desired)
	}
	//Apply rule where two equal sub exps on the left appear once on the right
	env = newenviron()
	_, _, rname = helper_undistribrule(env)
	ename = "ex"
	exp = "(+ (* a (+ c d)) (* a e))"
	i = "1"
	desired = "(* a (+ (+ c d) e))"
	resultname := "res"
	helper_defexps([]string{ename + " " + exp}, env)
	result = helper_applyrule(rname, ename, i, resultname, env)
	if desired != result {
		t.Errorf("The result of apply undistributive rule at %v is\n%s\n"+
			"It should be \n%s\n", i, result, desired)
	}
	if env.expmap[resultname] == nil {
		t.Error("Failed to store the result of a rule application")
	} else if env.expmap[resultname].String() != desired {
		t.Error("Rule result stored incorrectly")
	}
}
