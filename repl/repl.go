package repl

import (
	"bufio"
	"fmt"
	"io"
	"jeff/evaluator"
	"jeff/lexer"
	"jeff/object"
	"jeff/parser"
)

const PROMPT = ">>"

// Start the REPL. Keeps state so inputs can reuse variables
func Start(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(writer, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		lexer := lexer.New(line)
		parser := parser.New(lexer)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(writer, parser.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(writer, evaluated.Inspect())
			io.WriteString(writer, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
