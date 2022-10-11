package main

import (
	"fmt"
	"gomonkey/repl"
	"os"
	"os/user"
)

func main() {
	me, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is Monkey Programming だよ\n", me.Username)
	fmt.Printf("Feel free to type in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
