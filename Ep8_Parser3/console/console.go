package console

import (
	"bufio"
	"fmt"
	"interpreter/lexer"
	"interpreter/token"
	"io"
)

const PROMPT = ">> "

func Go(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		next := scanner.Scan()
		if !next {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.GetToken(); tok.Type != token.EOF; tok = l.GetToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
