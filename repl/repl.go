package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/rumpl/monkey-lang/lexer"
	"github.com/rumpl/monkey-lang/parser"
	"github.com/rumpl/monkey-lang/token"
)

const PROMPT = "🐒 "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		fmt.Fprintln(out, program.String())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintln(out, "\t"+msg)
	}
}
