package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

const PROMPT = "(pingul)>> "

func Start(in io.Reader, out io.Writer) {
	fmt.Printf("Welcome to the Pingul REPL!\n")
	fmt.Printf("Type 'exit' to quit the REPL.\n")

	for {
		fmt.Fprintf(out, PROMPT)

		scanner := bufio.NewScanner(in)
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			fmt.Fprintf(out, "Adios!\n")
			break
		}

		lxr := lexer.New(input)
		for tkn := lxr.NextToken(); tkn.Type != token.EOF; tkn = lxr.NextToken() {
			fmt.Fprintf(out, "%+v\n", tkn)
		}
	}
}
