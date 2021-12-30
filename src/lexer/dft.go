package lexer

import (
	"fmt"
)

var (
	ErrorInvalidInitialState    = fmt.Errorf("Provided an invalid initial state")
	ErrorInvalidFinalStateSet   = fmt.Errorf("Provided an invalid final state set")
	ErrorNotFinalState          = fmt.Errorf("Provided state does not represent a final state.")
	ErrorTransitionDoesNotExist = fmt.Errorf("Transition does not exist.")
)

func containsInt(slice []int, element int) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

type Dft struct {
	alphabet      []byte
	states        []int
	initialState  int
	finalStates   []int
	transitionMap map[int]map[byte]int
	currentState  int
}

func NewDft(alphabet []byte, states []int, initialState int, finalStates []int, transitionMap map[int]map[byte]int) (*Dft, error) {
	if !containsInt(states, initialState) {
		return &Dft{}, ErrorInvalidInitialState
	}

	if len(finalStates) > len(states) {
		return &Dft{}, ErrorInvalidFinalStateSet
	}

	return &Dft{
		alphabet:      alphabet,
		states:        states,
		initialState:  initialState,
		finalStates:   finalStates,
		transitionMap: transitionMap,
		currentState:  initialState,
	}, nil
}

func (d *Dft) transitionExists(char byte) bool {
	_, ok := d.transitionMap[d.currentState][char]
	return ok
}

// Next updates and returns the next state when consuming
// char in the current state. If there is no transitions
// possible to be made, Next returns the
// inital state and ErrorTransitionDoesNotExist
func (d *Dft) Next(char byte) (int, error) {
	if !d.transitionExists(char) {
		return d.initialState, ErrorTransitionDoesNotExist
	}
	newCurrentState := d.transitionMap[d.currentState][char]
	d.currentState = newCurrentState
	return d.currentState, nil
}

// Reset puts the dft on the inital state
func (d *Dft) Reset() {
	d.currentState = d.initialState
}

// IsFinalState returns whether we stopped on a final state or not
func (d *Dft) IsFinalState() bool {
	return containsInt(d.finalStates, d.currentState)
}

func (d *Dft) GetCurrentState() int {
	return d.currentState
}
