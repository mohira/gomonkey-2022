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

	fmt.Printf("Hello %s! REPLをどうぞ\n", me.Username)
	repl.Start(os.Stdin, os.Stdout)
}
