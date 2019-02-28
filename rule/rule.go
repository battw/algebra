package rule

import (
	"errors"
	"github.com/battw/algebra/expr"
	"strconv"
	"strings"
)

type Rule struct {
	lhs *expr.Expr
	rhs *expr.Expr
}

//TODO remove the error as it is never returned
func New(lhs, rhs *expr.Expr) (*Rule, error) {
	//Test for extra vars on rhs
	return &Rule{lhs, rhs}, nil
}

//Applicable - is the rule applicable to the designated sub expression of the given expression.
func (r *Rule) Applicable(e *expr.Expr, subi int) bool {
	mch, _ := r.lhs.Match(e.Subexp(subi))
	return mch
}

func (r *Rule) Apply(exp *expr.Expr, subi int, introexps ...*expr.Expr) (*expr.Expr, error) {
	//Associate variables in the rules lhs with subtrees of the sub expression
	match, varmap := r.lhs.Match(exp.Subexp(subi))
	if !match {
		return nil, errors.New("Rule\n" + r.String() +
			"\nis not applicable to \n" + exp.Subexp(subi).String())
	}
	var intros []rune = r.introductions()
	if len(intros) != len(introexps) {
		return nil, errors.New("Requires " + strconv.Itoa(len(intros)) + " new expressions." +
			"Got " + strconv.Itoa(len(introexps)))
	}
	for i, sym := range intros {
		varmap[sym] = introexps[i]
	}
	//Make a copy of the expression
	exp = exp.Clone()
	//Make a copy of rhs
	rhs := r.rhs.Clone()
	//Replace the variables in the copy of rhs with the associated subtrees
	rhs = rhs.Subvar(varmap)
	//Substitute the modified rhs back into the copied expression
	exp = exp.Substitute(subi, rhs)
	return exp, nil
}

func (r *Rule) String() string {
	return r.lhs.String() + " -> " + r.rhs.String()
}

func Pretty(r *Rule) string {
	lstr := expr.Pretty(r.lhs)
	rstr := expr.Pretty(r.rhs)
	llines := strings.Split(lstr, "\n")
	rlines := strings.Split(rstr, "\n")
	//remove empty lines
	for i := 0; i < len(llines); i++ {
		if len(llines[i]) == 0 {
			llines = append(llines[:i], llines[i+1:]...)
		}
	}
	for i := 0; i < len(rlines); i++ {
		if len(rlines[i]) == 0 {
			rlines = append(rlines[:i], rlines[i+1:]...)
		}
	}

	height := max(len(llines), len(rlines))
	diff := len(llines) - len(rlines)
	var small []string
	if diff < 0 {
		small = llines
	} else if diff > 0 {
		small = rlines
	}
	//Pad the smaller of the two expressions so that they are both the same height
	if small != nil {
		front := true // padding line at beginning insert
		for i := 0; i < abs(diff); i++ {
			line := strings.Repeat(" ", len(small[0]))
			front = !front
			if front {
				small = append([]string{line}, small...)
			} else {
				small = append(small, line)
			}
		}
	}
	if diff < 0 {
		llines = small
	} else if diff > 0 {
		rlines = small
	}
	//Create an arrow and the space around it to put between the two expressions
	arrow := make([]string, height)
	for i, _ := range arrow {
		if i == height/2 {
			arrow[i] = "  =>  "
		} else {
			arrow[i] = "      "
		}
	}
	pretty := ""
	for i := 0; i < height; i++ {
		pretty += llines[i] + arrow[i] + rlines[i] + "\n"
	}

	return pretty
}

//introsvar - If the rule introduces new subexpressions, return the variable symbol
//representing those subexpressions.
//Otherwise return an empty slice.
func (r *Rule) introductions() []rune {
	syms := make([]rune, 0)
	var lvars, rvars map[rune]bool = r.lhs.Vars(), r.rhs.Vars()
	for r, b := range rvars {
		if b && !lvars[r] {
			syms = append(syms, r)
		}
	}
	return syms
}

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
