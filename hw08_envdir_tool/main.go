package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		usage()
	}

	path := args[1]
	command := args[2]

	env, err := ReadDir(path)
	if err != nil {
		panic("Failed to read " + path + ": " + err.Error())
	}

	// Запустить программу
	// Связать stdin, stdout, stderr -> программа

	// make map from vars
	varsMap := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		varsMap[parts[0]] = strings.Join(parts[1:], "=")
	}

	// Process vars map
	for varName := range env {
		if env[varName].NeedRemove {
			delete(varsMap, varName)
		} else {
			varsMap[varName] = env[varName].Value
		}
	}

	// Convert map to []strings
	execVars := make([]string, 0)
	for varName := range varsMap {
		execVars = append(execVars, varName+"="+varsMap[varName])
	}

	cmd := exec.Command(command, args[3:]...)
	cmd.Env = execVars
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func usage() {
	fmt.Println("Usage: env_reader /env/dir")
	os.Exit(0)
}
