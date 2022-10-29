package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"gibbon/lexer"
	"gibbon/token"
	"io"
)

const PROMPT = ">> "

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	for {
		fmt.Fprint(output, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Bytes()
		l := lexer.NewLexer(bytes.NewReader(line), "REPL")
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(output, "%+v\n", tok)
		}
	}
}
