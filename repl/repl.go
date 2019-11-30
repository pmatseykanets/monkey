package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pmatseykanets/monkey/lexer"
	"github.com/pmatseykanets/monkey/token"
)

const PROMPT = ">> "

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

		lex := lexer.New(strings.NewReader(s.Text()))
		for {
			tok := lex.NextToken()
			if tok.Type == token.EOF {
				break
			}

			fmt.Fprintf(w, "%+v\n", tok)
		}
	}
}
