package console

import (
	"bufio"
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/parser"
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
		p := parser.New(l)

		program := p.ParseProgram()
		for _, element := range p.Errors() {
			fmt.Println(element)
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
