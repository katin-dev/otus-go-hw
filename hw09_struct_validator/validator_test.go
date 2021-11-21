package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			"Hello", // Простые типы не валидируем и не ругаемся на них
			nil,
		},
		{
			User{
				ID:     "A-10001",
				Name:   "Test Name",
				Age:    25,
				Email:  "katin.dev@gmail.com",
				Role:   "admin",
				Phones: []string{"79990001122"},
				meta:   []byte("{\"is_test\": true}"),
			},
			nil,
		},
		{
			User{
				ID:     "0000000000000000000000000000001111111", // 37 symbols
				Name:   "Test",
				Age:    16,                                      // < 18
				Email:  "some-invalid-email@",                   // invalid email
				Role:   "invalid_role",                          // not in enum
				Phones: []string{"79990001122", "+79990001122"}, // 2nd email is invalid: len > 11
				meta:   []byte("{\"is_test\": true}"),
			},
			ValidationErrors{
				ValidationError{"ID", ErrStrLen},
				ValidationError{"Age", ErrNumRange},
				ValidationError{"Email", ErrInvalidEmail},
				ValidationError{"Role", ErrStrEnum},
				ValidationError{"Phones.1", ErrStrLen},
			},
		},
		// @TODO добавить проверки на следующие структуры
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
