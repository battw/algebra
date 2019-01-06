package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	//    "strconv"
	"github.com/llo-oll/algebra/expr"
)

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

//readtok returns, as a string, the stream of runes up until the next unicode
//whitespace rune. The whitespace is discarded.
func readword(rch <-chan rune) string {
	var sb strings.Builder
	for r := range rch {
		if unicode.IsSpace(r) {
			break
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func readtok(rch <-chan rune) string {
	//ignore any initial whitespace
	var r rune
	for r = range rch {
		if !unicode.IsSpace(r) {
			break
		}
	}
	tok := ""

	//if starts with a ( then extract the expression string
	if r == '(' {
		tok += string(r)
		scope := 1
		for r := range rch {
			switch r {
			case '(':
				scope++
			case ')':
				scope--
			}
			tok += string(r)
			if scope == 0 {
				break
			}
		}
	} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
		tok += string(r)
		for r = range rch {
			if unicode.IsSpace(r) {
				break
			}
			tok += string(r)
		}
	}

	return tok
}

//tokenise, returns a channel providing tokens as strings.
//A token is either a space separated command/argument or an expr.
func tokenise(rch <-chan rune) <-chan string {
	tokch := make(chan string)
	go func() {
		for {
			tok := readtok(rch)
			if tok == "" {
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
