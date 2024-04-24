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
	fmt.Println("Type in commands")

	repl.Start(os.Stdin, os.Stdout)
	
}