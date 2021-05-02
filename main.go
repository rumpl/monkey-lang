package main

import (
	"fmt"
	"io/ioutil"

	"github.com/rumpl/monkey-lang/codegen"
	"github.com/rumpl/monkey-lang/lexer"
	"github.com/rumpl/monkey-lang/object"
	"github.com/rumpl/monkey-lang/parser"
)

func main() {
	compile("./monkey/hello.monkey")
	// fmt.Println("This is the Monkey programming language!")
	// repl.Start(os.Stdin, os.Stdout)
}

func compile(file string) {
	env := object.NewEnvironment()

	code, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	l := lexer.New(string(code))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	c := codegen.New(program)

	err = c.Codegen(env)
	if err != nil {
		fmt.Println(err)
	}
}

func printParserErrors(errors []string) {
	for _, msg := range errors {
		fmt.Println("\t" + msg)
	}
}
