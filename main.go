package main

import (
	"fmt"
	"os"

	"github.com/rumpl/monkey-lang/repl"
)

func main() {
	fmt.Println("This is the Monkey programming language!")
	repl.Start(os.Stdin, os.Stdout)
}
