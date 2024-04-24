package repl

import (
	"bufio"
	"fmt"
	"io"
	"jeff/lexer"
	"jeff/parser"
)

const PROMPT = ">>"

func Start(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

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

		io.WriteString(writer, program.String())
		io.WriteString(writer, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
