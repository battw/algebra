package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"
    "strconv"
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

func tokenise(input string) []string {
    //find index of first ( if present
    input = strings.Trim(input, "\n")
    expstart := strings.Index(input, " (") //CAREFUL the space is needed
    if expstart == -1 {
        return strings.Split(input, " ")
    }
    toks := strings.Split(input[:expstart], " ")
    //TODO filter tokens to remove empty strings
    exp := input[expstart:]
    return append(toks, exp)

}

func handle(input string, expmap map[string]*expr.Expr) string {
    toks := tokenise(input)
    if len(toks) == 0 {
        return ""
    }
    switch toks[0] {
    case "exp", "e":
        if len(toks) < 2 {
            return "exp must be followed by an expression"
        }
        expmap["exp"] = expr.Translate(toks[1])
        return expmap["exp"].String()
    case "sub", "s":
        //TODO sort this monstrosity of a case
        if len(toks) < 2 {
            return "sub must be followed by an integer"
        }
        index, err := strconv.Atoi(toks[1])
        if err != nil {
            return "sub must be followed by an integer"
        }
        //TODO handle the possible nil returned by Sub
        return expmap["exp"].Sub(index).String()
    }

    return ""
}

func main() {
	engin, engout := eng()
	repl(engin, engout)
}
