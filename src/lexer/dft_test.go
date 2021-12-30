package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func prepareDft() *Dft {
	return &Dft{
		alphabet:     []byte{'A', 'B', 'C', 'D', 'F', 'G'},
		states:       []int{0, 1, 2, 3},
		initialState: 0,
		finalStates:  []int{2, 3},
		transitionMap: map[int]map[byte]int{
			0: {
				'A': 1,
				'B': 2,
				'C': 1,
				'D': 1,
			},
			1: {
				'A': 1,
				'B': 2,
			},
			2: {
				'B': 2,
				'C': 3,
			},
			3: {
				'A': 3,
				'B': 3,
				'C': 3,
				'D': 3,
				'F': 3,
				'G': 3,
			},
		},
	}
}

func TestDft_Next(t *testing.T) {
	dft := prepareDft()

	testCases := []struct {
		name          string
		charsToRead   []byte
		expectedError error
		expectedState int
		errorIndex    int
	}{
		{
			name:          "Next transition valid with one character",
			charsToRead:   []byte{'A'},
			expectedError: nil,
			expectedState: 1,
			errorIndex:    -1,
		},
		{
			name:          "Next transition invalid",
			charsToRead:   []byte{'L'},
			expectedError: ErrorTransitionDoesNotExist,
			expectedState: dft.initialState,
			errorIndex:    0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := require.New(t)

			var currentState int
			var err error

			for idx, c := range testCase.charsToRead {
				currentState, err = dft.Next(c)
				if idx == testCase.errorIndex {
					r.Equal(testCase.expectedError, err)
				}
			}

			r.Equal(testCase.expectedState, currentState)
			dft.currentState = dft.initialState
		})
	}
}
