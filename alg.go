package main

import (
	"bufio"
	"fmt"
	"os"
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
	var exp *expr.Expr
	go func() {
	    for instr := range engin {
			exp = expr.Translate(instr)
			engout <- exp.String()
		}
	}()
	return engin, engout
}

func main() {
	engin, engout := eng()
	repl(engin, engout)
}
