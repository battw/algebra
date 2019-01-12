package main

import (
	"bufio"
	"fmt"
	"github.com/llo-oll/algebra/expr"
	"github.com/llo-oll/algebra/toknify"
	"os"
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

func handle(input string, expmap map[string]*expr.Expr) string {
	runech := expr.Strstream(input)
	tokch := toknify.Tokenise(runech)
	for tok := range tokch {
		fmt.Println(tok)
	}
	return " "
}

func main() {
	engin, engout := eng()
	repl(engin, engout)
}
