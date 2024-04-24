package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
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
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	
	fmt.Print(REPL_HEADER)
	fmt.Printf("Hello %s, Welcome to the Jeff programming language!\n", user.Username)
	fmt.Println("Type in commands, Type 'exit' to close")

	repl.Start(os.Stdin, os.Stdout)
	
}