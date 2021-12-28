package lexer

import (
	"fmt"
)

var (
	ErrorNotFinalState          = fmt.Errorf("Provided state does not represent a final state.")
	ErrorTransitionDoesNotExist = fmt.Errorf("Transition does not exist.")
)

var dft = map[int]map[byte]int{
	0: {
		'\n': 0,
		'\t': 0,
		' ':  0,
		'L':  1, // L is the set that represents all the letters of the alphabet, low and capital case
	},
	1: {
		'L': 1,
		'D': 1, // D is the set that represents all numbers from 0 to 9
		'_': 1,
	},
}

var finalStateTable = map[int]TokenClass{
	1: IDENTIFIER,
}

// Get the token class of a given final state
func GetTokenClass(state int) (TokenClass, error) {
	tokenClass, exists := finalStateTable[state]

	if !exists {
		return ERROR, ErrorNotFinalState
	}

	return tokenClass, nil
}

// Get the next state given the current state and symbol
func GetNextState(currState int, symbol byte) (int, error) {
	nextState, exists := dft[currState][symbol]

	if !exists {
		return currState, ErrorTransitionDoesNotExist
	}

	return nextState, nil
}
