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
		Version string   `validate:"regexp:\\d+|len:5"`
		Phones  []string `validate:"len:5"`
		Header  []int
		Code    int `validate:"min:0|max:10"`
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
			in: App{
				Version: "123456",
				Phones:  []string{"s12333", "122268"},
				Header:  []int{1, 2, 3, 10},
				Code:    17,
			},
			expectedErr: fmt.Errorf(
				"field: Version |||| err: could not len tagStringValidate : rule for tag validate not match\n" +
					"field: Phones |||| err: could not len tagStringValidate : rule for tag validate not match\n" +
					"field: Phones |||| err: could not len tagStringValidate : rule for tag validate not match\n" +
					"field: Code |||| err: could not max tagIntValidate : rule for tag validate not match\n"),
		},
		{
			in: App{
				Version: "12345",
				Phones:  []string{"12333", "12268"},
				Header:  []int{1, 2, 3, 10},
				Code:    7,
			},
			expectedErr: fmt.Errorf(""),
		},
		{
			in: Response{
				Body: "sdasddsa",
				Code: 7,
			},
			expectedErr: fmt.Errorf(
				"field: Code |||| err: could not in tagIntValidate :rule for tag not valid\n"), //nolint: revive
		},
		{
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{1, 2, 3},
				Signature: []byte{1, 2, 3},
			},
			expectedErr: fmt.Errorf(""),
		},
		{
			in: User{
				ID:     "not36len",
				Name:   "vlad",
				Age:    21,
				Email:  "arsenal@gmail.ru",
				Role:   "admin",
				Phones: []string{"12345", "12333", "8797"},
				meta:   nil,
			},
			expectedErr: fmt.Errorf(
				"field: ID |||| err: could not len tagStringValidate : rule for tag validate not match\n" +
					"field: Phones |||| err: could not len tagStringValidate : rule for tag validate not match\n" +
					"field: Phones |||| err: could not len tagStringValidate : rule for tag validate not match\n" +
					"field: Phones |||| err: could not len tagStringValidate : rule for tag validate not match\n"),
		},
		{
			in: User{
				ID:     "not36lenqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
				Name:   "vlad",
				Age:    21,
				Email:  "arsenal@gmail.ru",
				Role:   "admin",
				Phones: []string{"12345678911", "12333456329", "87971234000"},
				meta:   nil,
			},
			expectedErr: fmt.Errorf(""),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.EqualError(t, Validate(tt.in), tt.expectedErr.Error())
		})
	}
}
