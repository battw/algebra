package toknify

import (
	"fmt"
	"strconv"
	//	"strings"
	"unicode"
)

type Toktyp int

const (
	NIL Toktyp = iota
	ERR
	EXPR
	INT
	NAME
)

func (tt Toktyp) String() string {
	switch tt {
	case NIL:
		return "NIL"
	case ERR:
		return "ERR"
	case EXPR:
		return "EXPR" //starts with a (
	case INT:
		return "INT"
	case NAME:
		return "NAME" //starts with a letter
	default:
		return "INVALID TYPE"
	}
}

type Tokn struct {
	Typ Toktyp
	Str string
}

func (tok Tokn) String() string {
	return fmt.Sprintf("Tokn{%s, %s}", tok.Typ, tok.Str)
}

//tokenise, returns a channel providing tokens as strings.
//A token is either a space separated command/argument or an expr.
func Tokenise(rch <-chan rune) <-chan Tokn {
	tokch := make(chan Tokn)
	go func() {
		for {
			tokstr := readtokstr(rch)
			tok := str2Tokn(tokstr)
			if tok.Typ == NIL {
				break
			}
			tokch <- tok
		}
		close(tokch)
	}()
	return tokch
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

func str2Tokn(str string) Tokn {
	var tok Tokn
	tok.Str = str
	if len(str) == 0 {
		tok.Typ = NIL
	} else if isInt(str) {
		tok.Typ = INT
	} else if isName(str) {
		tok.Typ = NAME
	} else if isExpr(str) {
		tok.Typ = EXPR
	} else {
		tok.Typ = ERR
		tok.Str = "Token string '" + str + "' isn't a NAME, INT or EXPR"
	}
	return tok
}
