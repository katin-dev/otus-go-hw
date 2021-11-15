package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
	}

	path := args[1]

	env, err := ReadDir(path)
	if err != nil {
		panic("Failed to read " + path + ": " + err.Error())
	}

	fmt.Println(env)
}

func usage() {
	fmt.Println("Usage: env_reader /env/dir")
	os.Exit(0)
}
