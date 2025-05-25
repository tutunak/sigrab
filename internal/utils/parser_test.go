package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIssueKey(t *testing.T) {
	tests := []struct {
		name           string
		issueKey       string
		expectedError  bool
		expectedPrefix string
		expectedNumber int
	}{
		// Success cases
		{
			name:           "standard issue key",
			issueKey:       "ABC-123",
			expectedError:  false,
			expectedPrefix: "ABC",
			expectedNumber: 123,
		},
		{
			name:           "single digit issue number",
			issueKey:       "DEV-1",
			expectedError:  false,
			expectedPrefix: "DEV",
			expectedNumber: 1,
		},
		{
			name:           "large issue number",
			issueKey:       "PROJ-9999",
			expectedError:  false,
			expectedPrefix: "PROJ",
			expectedNumber: 9999,
		},
		{
			name:           "lowercase prefix",
			issueKey:       "test-123",
			expectedError:  false,
			expectedPrefix: "test",
			expectedNumber: 123,
		},

		// Error cases
		{
			name:          "empty string",
			issueKey:      "",
			expectedError: true,
		},
		{
			name:          "missing hyphen",
			issueKey:      "ABC123",
			expectedError: true,
		},
		{
			name:          "multiple hyphens",
			issueKey:      "ABC-123-XYZ",
			expectedError: true,
		},
		{
			name:          "missing prefix",
			issueKey:      "-123",
			expectedError: true,
		},
		{
			name:          "missing number",
			issueKey:      "ABC-",
			expectedError: true,
		},
		{
			name:          "non-numeric part after hyphen",
			issueKey:      "ABC-XYZ",
			expectedError: true,
		},
		{
			name:          "negative number",
			issueKey:      "ABC--123",
			expectedError: true,
		},
		{
			name:          "decimal number",
			issueKey:      "ABC-123.45",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, number, err := ParseIssueKey(tt.issueKey)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPrefix, prefix)
				assert.Equal(t, tt.expectedNumber, number)
			}
		})
	}
}
