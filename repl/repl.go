package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pmatseykanets/monkey/eval"
	"github.com/pmatseykanets/monkey/lexer"
	"github.com/pmatseykanets/monkey/parser"
)

const PROMPT = ">> "

// Start .
func Start(r io.Reader, w io.Writer) {
	s := bufio.NewScanner(r)

	for {
		fmt.Fprint(w, PROMPT)
		if !s.Scan() {
			return
		}
		if s.Text() == "\\q" {
			return
		}

		p := parser.New(lexer.FromString(s.Text()))
		prg := p.Parse()
		if len(p.Errors()) > 0 {
			fmt.Fprintln(w, "parser errors:")
			for _, msg := range p.Errors() {
				fmt.Fprintln(w, "\t"+msg.Error())
			}
			continue
		}

		evald := eval.Eval(prg)
		if evald != nil {
			fmt.Fprintln(w, evald.Inspect())
		}
	}
}
