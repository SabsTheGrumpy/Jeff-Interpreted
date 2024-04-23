package repl

import (
	"bufio"
	"fmt"
	"monkey/token"
	"io"
	"monkey/lexer"
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
		lexer := lexer.New(line)

		for tok := lexer.NextToken(); tok.Type != token.EOF; tok = lexer.NextToken() {
			fmt.Fprintf(writer, "%+v\n", tok)
		}
	}

}