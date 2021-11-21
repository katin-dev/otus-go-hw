package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrStrLen       error = errors.New("string length does not match expected")
	ErrNumRange     error = errors.New("number is out of range")
	ErrInvalidEmail error = errors.New("invalid email")
	ErrStrEnum      error = errors.New("the value is not in allowed enum")
	ErrRegexp       error = errors.New("the value does not match regexp pattern")
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
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Err.Error())
	}
	return strings.Join(messages, "\n")
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
		v := vValue.Field(i)
		tag := f.Tag.Get("validate")

		if tag == "" {
			continue
		}

		if !isAllowedType(f.Type) {
			continue
		}

		validators, err := parseTag(tag)
		if err != nil {
			return err
		}

		for _, validator := range validators {
			if f.Type.Kind() == reflect.Slice {
				for j := 0; j < vValue.Field(i).Len(); j++ {
					validateField(f.Name+"."+strconv.Itoa(j), vValue.Field(i).Index(j), validator, &errors)
				}
			} else {
				validateField(f.Name, v, validator, &errors)
			}
		}

	}

	if errors.Length() == 0 {
		return nil
	}

	return errors
}

func validateField(name string, v reflect.Value, validator Validator, errors *ValidationErrors) {
	validationErr := validator.Validate(v)
	if validationErr != nil {
		if validationErrors, ok := validationErr.(ValidationErrors); ok {
			*errors = append(*errors, validationErrors...)
		} else {
			*errors = append(*errors, ValidationError{name, validationErr})
		}
	}
}

func isAllowedType(t reflect.Type) bool {
	return t.Kind() == reflect.Int ||
		t.Kind() == reflect.String ||
		t == reflect.SliceOf(reflect.TypeOf("")) ||
		t == reflect.SliceOf(reflect.TypeOf(1)) ||
		t.Kind() == reflect.Struct
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
		if len(optionParts) > 1 {
			validatorOptionString = optionParts[1]
		}

		var validator Validator
		var err error

		switch validatorName {
		case "len":
			validator, err = NewStrLenValidator(validatorOptionString)
		case "min":
			validator, err = NewNumMinValidator(validatorOptionString)
		case "max":
			validator, err = NewNumMaxValidator(validatorOptionString)
		case "regexp":
			validator, err = NewRegExpValidator(validatorOptionString)
		case "in":
			validator, err = NewStrEnumValidator(validatorOptionString)
		case "nested":
			validator, err = NewNestedValidator(validatorOptionString)
		default:
			validator = nil
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create validator %s: %s", validatorName, err)
		}

		if validator != nil {
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

	if length, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &StrLenValidator{length}, nil
}

func (v *StrLenValidator) Validate(val reflect.Value) error {
	strVal := val.String()
	if len(strVal) != v.len {
		return ErrStrLen
	}

	return nil
}

type NumMinvalidator struct {
	min int64
}

func NewNumMinValidator(options string) (*NumMinvalidator, error) {
	options = strings.TrimSpace(options)

	var (
		num int
		err error
	)

	if num, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &NumMinvalidator{int64(num)}, nil
}

func (v *NumMinvalidator) Validate(val reflect.Value) error {
	intVal := val.Int()
	if intVal < v.min {
		return ErrNumRange
	}

	return nil
}

type NumMaxValidator struct {
	max int64
}

func NewNumMaxValidator(options string) (*NumMaxValidator, error) {
	options = strings.TrimSpace(options)

	var (
		num int
		err error
	)

	if num, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &NumMaxValidator{int64(num)}, nil
}

func (v *NumMaxValidator) Validate(val reflect.Value) error {
	intVal := val.Int()
	if intVal > v.max {
		return ErrNumRange
	}

	return nil
}

type RegExpValidator struct {
	re *regexp.Regexp
}

func NewRegExpValidator(options string) (*RegExpValidator, error) {
	options = strings.TrimSpace(options)

	re, err := regexp.Compile(options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regexp: %s", err)
	}

	return &RegExpValidator{re}, nil
}

func (v *RegExpValidator) Validate(val reflect.Value) error {
	value := val.String()

	if !v.re.Match([]byte(value)) {
		return ErrRegexp
	}

	return nil
}

type StrEnumValidator struct {
	enums []string
}

func NewStrEnumValidator(options string) (*StrEnumValidator, error) {
	options = strings.TrimSpace(options)

	return &StrEnumValidator{strings.Split(options, ",")}, nil
}

func (v *StrEnumValidator) Validate(val reflect.Value) error {
	var value string

	if val.Kind() == reflect.Int {
		value = strconv.Itoa(int(val.Int()))
	} else {
		// Для простоты не будем проверять на другие типы
		value = val.String()
	}

	for _, expected := range v.enums {
		if value == expected {
			return nil
		}
	}

	return ErrStrEnum
}

type NestedValidator struct{}

func NewNestedValidator(options string) (*NestedValidator, error) {
	// options = strings.TrimSpace(options)

	return &NestedValidator{}, nil
}

func (v *NestedValidator) Validate(val reflect.Value) error {
	return Validate(val.Interface())
}
