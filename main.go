package main

import (
	"fmt"
	"jeff/evaluator"
	"jeff/lexer"
	"jeff/object"
	"jeff/parser"
	"jeff/repl"
	"os"
	"os/user"
	"strings"
)

const REPL_HEADER = `
  |||   ||||||||   ||||||||   ||||||||
  |||   |||        |||        |||
  |||   ||||||||   ||||||||   ||||||||
 ///    |||        |||        |||
///     ||||||||   |||        |||
`

// Simple Repl
func main() {

	if len(os.Args) < 2 {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Print(REPL_HEADER)
		fmt.Printf("Hello %s, Welcome to the Jeff programming language!\n", user.Username)
		fmt.Println("Type in commands, Type 'exit' to close")

		repl.Start(os.Stdin, os.Stdout)
	} else if len(os.Args) == 2 {
		fileName := os.Args[1]
		if !strings.HasSuffix(fileName, ".jeff") {
			fmt.Printf("ERROR: file %s is not a .jeff file\n", fileName)
			return
		}

		env := object.NewEnvironment()
		data, err := os.ReadFile(fileName)

		if err != nil {
			fmt.Printf("ERROR: file %s can't be read\n", fileName)
			return
		}

		lexer := lexer.New(string(data))
		parser := parser.New(lexer)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			repl.PrintParserErrors(os.Stdout, parser.Errors())
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			// io.WriteString(writer, evaluated.Inspect())
			// io.WriteString(writer, "\n")
			fmt.Print(evaluated.Inspect())
		}
	}
}
