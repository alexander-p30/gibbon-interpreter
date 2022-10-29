package main

import (
	"fmt"
	"gibbon/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to the gibbon interpreter %s!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
