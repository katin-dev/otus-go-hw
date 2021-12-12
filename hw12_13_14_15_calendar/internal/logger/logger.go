package logger

import "fmt"

type Logger struct {
	file  string
	level string
}

func New(file, level string) *Logger {
	return &Logger{
		file:  file,
		level: level,
	}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
