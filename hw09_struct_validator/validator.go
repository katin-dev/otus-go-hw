package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	vValue := reflect.ValueOf(v)

	// Валидировать будем только структуры
	if vType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < vType.NumField(); i++ {
		f := vType.Field(i)
		tag := f.Tag.Get("validate")

		if tag == "" {
			continue
		}

		validators, err := parseTag(tag)
		if err != nil {
			return err
		}

		for _, validator := range validators {
			if f.Type.Kind() == reflect.Slice {
				for j := 0; j < vValue.Field(i).Len(); j++ {
					validationErr := validator.Validate(vValue.Field(i).Index(j))
					if validationErr != nil {
						errors = append(errors, ValidationError{f.Name, validationErr})
					}
				}
			} else {
				validationErr := validator.Validate(vValue.Field(i))
				if validationErr != nil {
					errors = append(errors, ValidationError{f.Name, validationErr})
				}
			}
		}

	}

	if errors.Length() == 0 {
		return nil
	}

	return errors
}

func parseTag(tag string) ([]Validator, error) {
	validators := make([]Validator, 0)

	parts := strings.Split(tag, "|")
	for _, part := range parts {

		optionParts := strings.Split(part, ":")

		if len(optionParts) == 0 {
			return nil, fmt.Errorf("failed to parse validator: " + part)
		}

		validatorName := optionParts[0]

		validatorOptionString := ""
		if len(optionParts) > 0 {
			validatorOptionString = optionParts[1]
		}

		switch validatorName {
		case "len":
			validator, err := NewStrLenValidator(validatorOptionString)
			if err != nil {
				return nil, fmt.Errorf("failed to create validator %s: %s", validatorName, err)
			}

			validators = append(validators, validator)
		}
	}

	return validators, nil
}

type Validator interface {
	Validate(v reflect.Value) error
}

type StrLenValidator struct {
	len int
}

func NewStrLenValidator(options string) (*StrLenValidator, error) {
	options = strings.TrimSpace(options)

	var (
		length int
		err    error
	)

	fmt.Println("Option: ", options)

	if length, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &StrLenValidator{length}, nil
}

func (v *StrLenValidator) Validate(val reflect.Value) error {
	strVal := val.String()
	fmt.Println("StrLenValidator: ", strVal, v.len)
	if len(strVal) > v.len {
		return ErrStrLen
	}

	return nil
}

var (
	ErrStrLen       error = errors.New("string length exceed the limit")
	ErrNumRange     error = errors.New("number is out of range")
	ErrInvalidEmail error = errors.New("invalid email")
	ErrStrEnum      error = errors.New("the value is not in allowed enum")
)
