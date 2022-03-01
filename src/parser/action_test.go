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
		lexem           string
		expectedAction  Action
		expectedOperand int
	}{
		{
			name:            "Getting accept",
			state:           1,
			lexem:           "$",
			expectedAction:  ACCEPT,
			expectedOperand: 0,
		},
		{
			name:            "Get shift",
			state:           3,
			lexem:           "leia",
			expectedAction:  SHIFT,
			expectedOperand: 11,
		},
		{
			name:            "Get reduce",
			state:           23,
			lexem:           "id",
			expectedAction:  REDUCE,
			expectedOperand: 8,
		},
		{
			name:            "Get none",
			state:           3,
			lexem:           "varinicio",
			expectedAction:  NONE,
			expectedOperand: 0,
		},
	}

	action := NewActionReader(actionTablePath)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			action, opr := action.GetAction(lexer.State(tc.state), lexer.NewToken(lexer.NUM, tc.lexem, lexer.NULL))
			r.Equal(tc.expectedAction, action)
			r.Equal(tc.expectedOperand, opr)
		})
	}
}
