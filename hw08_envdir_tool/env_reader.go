package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		val := EnvValue{}
		if file.Size() == 0 {
			val.NeedRemove = true
			env[file.Name()] = val
			continue
		}

		content, _ := os.ReadFile(dir + string(os.PathSeparator) + file.Name())
		val.Value = string(content)
		val.Value = strings.Split(val.Value, "\n")[0]
		val.Value = strings.TrimRight(val.Value, " ")
		val.Value = string(bytes.Replace([]byte(val.Value), []byte{0x00}, []byte{'\n'}, -1))
		env[file.Name()] = val
	}

	return env, nil
}
