package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateMessage(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedOutput error
	}{
		{
			name:           "Should return an error when user id is empty",
			input:          "",
			expectedOutput: fmt.Errorf("userID can't be empty"),
		},
		{
			name:           "Should return a nil error when user id not empty",
			input:          "valid user id",
			expectedOutput: nil,
		},
	}

	messageValidator := NewMessageValidatorImpl()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Operation
			output := messageValidator.ValidateMessage(c.input)

			// Validation
			assert.EqualValues(t, c.expectedOutput, output)
		})
	}
}
