package main

import (
	"fmt"
	"os"

	"github.com/aziflaj/pingul/eval"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/object"
	"github.com/aziflaj/pingul/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		return
	}

	filename := os.Args[1]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist!\n", filename)
		return
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	lxr := lexer.New(string(content))
	p := parser.New(lxr)
	program := p.ParseProgram()

	scope := object.NewScope()
	result := eval.Eval(scope, program)

	if result != nil {
		fmt.Println(result.Inspect())
	} else {
		fmt.Println("NADA")
	}

	fmt.Println("Done!")
}
