package main

import (
	"os"

	"github.com/aziflaj/pingul/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
