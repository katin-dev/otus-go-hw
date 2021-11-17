package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		usage()
	}

	env, err := ReadDir(args[1])
	if err != nil {
		panic("Failed to read '" + args[1] + "': " + err.Error())
	}

	code := RunCmd(args[2:], env)

	os.Exit(code)
}

func usage() {
	fmt.Println("Usage: go-envdir /path/to/env/dir /path/to/command command arguments")
	os.Exit(0)
}
