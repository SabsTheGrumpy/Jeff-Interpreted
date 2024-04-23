package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)


// Simple Repl
func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s, Welcome to the Jeff programming language;\n", user.Username)
	fmt.Println("Type in commands")

	repl.Start(os.Stdin, os.Stdout)
	
}