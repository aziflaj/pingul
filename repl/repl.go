package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/aziflaj/pingul/eval"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/object"
	"github.com/aziflaj/pingul/parser"
)

const PROMPT = "(pingul)>> "

func Start(in io.Reader, out io.Writer) {
	fmt.Printf("Welcome to the Pingul REPL!\n")
	fmt.Printf("Type 'exit' to quit the REPL.\n")

	for {
		fmt.Fprintf(out, PROMPT)

		scanner := bufio.NewScanner(in)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(out, "Error reading input: %s\n", err)
			} else {
				fmt.Fprintf(out, "\nAdios!\n")
			}
			break
		}

		input := scanner.Text()

		if input == "exit" {
			fmt.Fprintf(out, "Adios!\n")
			break
		}

		lxr := lexer.New(input)
		p := parser.New(lxr)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		globalScope := object.NewScope()
		result := eval.Eval(globalScope, program)

		fmt.Fprintf(out, result.Inspect())
		fmt.Fprintf(out, "\n\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
