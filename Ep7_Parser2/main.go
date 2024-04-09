package main

import (
	"fmt"
	"interpreter/console"
	"os"
	"os/user"
)

func main() {
	_, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Type in commands\n")
	console.Go(os.Stdin, os.Stdout)
}
