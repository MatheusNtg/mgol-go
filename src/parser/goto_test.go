package parser

import (
	"mgol-go/src/lexer"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	gotoTablePath = "./tables/goto.tsv"
)

func TestGetGoto(t *testing.T) {
	testCases := []struct {
		name          string
		inicialState  int
		nonTerminal   string
		expectedState int
	}{
		{
			name:          "Getting Valid State",
			inicialState:  0,
			nonTerminal:   "P",
			expectedState: 1,
		},
		{
			name:          "Getting Valid State 2",
			inicialState:  21,
			nonTerminal:   "L",
			expectedState: 49,
		},
		{
			name:          "Getting Non Existent State",
			inicialState:  0,
			nonTerminal:   "V",
			expectedState: -1,
		},
	}

	got := NewGotoReader(gotoTablePath)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			ret_state := got.GetGoto(lexer.State(tc.inicialState), tc.nonTerminal)
			r.Equal(tc.expectedState, ret_state)
		})
	}
}
