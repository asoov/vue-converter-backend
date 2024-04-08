package helpers

import (
	"testing"
)

func TestCalculateNeededTokens(t *testing.T) {
	type test struct {
		naming            string
		input             string
		expected          int
		calculateFunction func(str string) (int, error)
		expectErr         bool
	}
	tests := []test{
		{
			naming:   "Successful",
			input:    "Lorem Ipsum Dolor sit Amet",
			expected: 10,
			calculateFunction: func(str string) (int, error) {
				return 10, nil
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.naming, func(t *testing.T) {
			tokens, err := CalculateNeededTokens(tc.input, tc.calculateFunction)
			if tc.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tc.expectErr && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			} else if tokens != tc.expected {
				t.Errorf("Expected %d tokens, got %d", tc.expected, tokens)
			}
		})
	}

}
