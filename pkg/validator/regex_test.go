package validator_test

import (
	"testing"
	"transfer-system/pkg/validator"
)

func TestValidateDecimalFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid cases
		{
			name:     "Valid format - integer part only one digit",
			input:    "1.12345",
			expected: true,
		},
		{
			name:     "Valid format - integer part multiple digits",
			input:    "12345.12345",
			expected: true,
		},
		{
			name:     "Valid format - large integer part",
			input:    "987654321.12345",
			expected: true,
		},
		{
			name:     "Valid format - zero integer part",
			input:    "0.12345",
			expected: true,
		},

		// Invalid cases
		{
			name:     "Invalid - less than 5 decimal digits",
			input:    "123.1234",
			expected: false,
		},
		{
			name:     "Invalid - more than 5 decimal digits",
			input:    "123.123456",
			expected: false,
		},
		{
			name:     "Invalid - no decimal part",
			input:    "123",
			expected: false,
		},
		{
			name:     "Invalid - no integer part",
			input:    ".12345",
			expected: false,
		},
		{
			name:     "Invalid - comma instead of dot",
			input:    "123,12345",
			expected: false,
		},
		{
			name:     "Invalid - includes letters",
			input:    "abc.12345",
			expected: false,
		},
		{
			name:     "Invalid - empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid - negative number",
			input:    "-123.12345",
			expected: false,
		},
		{
			name:     "Invalid - contains special characters",
			input:    "123!.12345",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.ValidateDecimalFormat(tt.input)

			if got != tt.expected {
				t.Errorf("ValidateDecimalFormat(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
