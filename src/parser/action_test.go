package parser

import (
	"mgol-go/src/lexer"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	actionTablePath = "./tables/action.tsv"
)

func TestGetAction(t *testing.T) {
	testCases := []struct {
		name            string
		state           int
		tokenClass      lexer.TokenClass
		expectedAction  Action
		expectedOperand int
	}{
		{
			name:            "Getting accept",
			state:           1,
			tokenClass:      lexer.EOF,
			expectedAction:  ACCEPT,
			expectedOperand: 0,
		},
		{
			name:            "Get shift",
			state:           3,
			tokenClass:      "leia",
			expectedAction:  SHIFT,
			expectedOperand: 11,
		},
		{
			name:            "Get reduce",
			state:           23,
			tokenClass:      lexer.IDENTIFIER,
			expectedAction:  REDUCE,
			expectedOperand: 8,
		},
		{
			name:            "Get none",
			state:           3,
			tokenClass:      "varinicio",
			expectedAction:  NONE,
			expectedOperand: 0,
		},
	}

	action := NewActionReader(actionTablePath)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			action, opr := action.GetAction(lexer.State(tc.state), lexer.NewToken(tc.tokenClass, "", lexer.NULL))
			r.Equal(tc.expectedAction, action)
			r.Equal(tc.expectedOperand, opr)
		})
	}
}
