package errorhandling

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetErrorType(t *testing.T) {
	testCases := []struct {
		name         string
		lexem        string
		expectedType LexicalErrorType
	}{
		{
			name:         "Lexical error",
			lexem:        `"this is an error`,
			expectedType: InvalidLiteral,
		},
		{
			name:         "Default number error",
			lexem:        "123123.",
			expectedType: InvalidNumber,
		},
		{
			name:         "Number error with letter",
			lexem:        "1231e",
			expectedType: InvalidNumber,
		},
		{
			name:         "Invalid comment",
			lexem:        "{asdfasdf",
			expectedType: InvalidComment,
		},
		{
			name:         "Invalid word",
			lexem:        "adaweqw$",
			expectedType: InvalidWord,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			actualValue := getErrorType(tc.lexem)
			r.Equal(tc.expectedType, actualValue)
		})
	}
}
