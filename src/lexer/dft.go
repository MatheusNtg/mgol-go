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

type State int
type Symbol byte

type Transition struct {
	from    State
	to      State
	reading []Symbol
}

type Dft struct {
	alphabet      []Symbol
	states        []State
	initialState  State
	finalStates   []State
	transitionMap map[State][]Transition
	currentState  State
}

func NewDft(alphabet []Symbol, states []State, initialState State, finalStates []State, transitionMap map[State][]Transition) (*Dft, error) {
	if !ContainsState(states, initialState) {
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

// Checks the existence of a certain symbol
// inside a reading slice in a Transition
func (d *Dft) transitionExists(char Symbol) bool {
	possibleTransitions := d.transitionMap[d.currentState]
	for _, transition := range possibleTransitions {
		for _, symbol := range transition.reading {
			if char == symbol {
				return true
			}
		}
	}
	return false
}

func (d *Dft) getNextTransition(char Symbol) Transition {
	possibleTransitions := d.transitionMap[d.currentState]
	for _, transition := range possibleTransitions {
		for _, symbol := range transition.reading {
			if char == symbol {
				return transition
			}
		}
	}
	return Transition{}
}

// Next updates and returns the next state when consuming
// char in the current state. If there is no transitions
// possible to be made, Next returns the
// inital state and ErrorTransitionDoesNotExist
func (d *Dft) Next(char Symbol) (State, error) {
	if !d.transitionExists(char) {
		return d.initialState, ErrorTransitionDoesNotExist
	}

	nextTransition := d.getNextTransition(char)
	d.currentState = nextTransition.to

	return d.currentState, nil
}

// Reset puts the dft on the inital state
func (d *Dft) Reset() {
	d.currentState = d.initialState
}

// IsFinalState returns whether we stopped on a final state or not
func (d *Dft) IsFinalState() bool {
	return ContainsState(d.finalStates, d.currentState)
}

func (d *Dft) GetCurrentState() State {
	return d.currentState
}
