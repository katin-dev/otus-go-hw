package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Length() int {
	return len(v)
}

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	errors := make(ValidationErrors, 0)

	vType := reflect.TypeOf(v)

	// Валидировать будем только структуры
	if vType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < vType.NumField(); i++ {
		f := vType.Field(i)
		tag := f.Tag
		fmt.Println(f.Name + " " + string(tag))
	}

	if errors.Length() == 0 {
		return nil
	}

	return errors
}

var (
	ErrStrLen       error = errors.New("string length exceed the limit")
	ErrNumRange     error = errors.New("number is out of range")
	ErrInvalidEmail error = errors.New("invalid email")
	ErrStrEnum      error = errors.New("the value is not in allowed enum")
)
