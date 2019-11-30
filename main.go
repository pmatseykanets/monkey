package main

import (
	"fmt"
	"os"

	"github.com/pmatseykanets/monkey/repl"
)

func main() {
	fmt.Fprint(os.Stdout, "Monkey REPL\n")
	repl.Start(os.Stdin, os.Stdout)
}
