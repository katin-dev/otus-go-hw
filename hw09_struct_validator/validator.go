package hw09structvalidator

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	// 1. Прочитать все поля структуры v (а если там не структура?)
	return ValidationErrors{
		ValidationError{"name", fmt.Errorf("Invalid")},
	}
}

var (
	ErrStrLen       error = errors.New("string length exceed the limit")
	ErrNumRange     error = errors.New("number is out of range")
	ErrInvalidEmail error = errors.New("invalid email")
	ErrStrEnum      error = errors.New("the value is not in allowed enum")
)
