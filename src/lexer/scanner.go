package lexer

import (
	"io"
	"log"
	"os"
)

// Recognize if a symbol is in the Letter set or the Digit set
//
// The Letter set is defined as L: {a, b, ..., z, A, B, ..., Z}
//
// The Digit set is defined as D: {0, 1, ..., 9}
func RecognizeSymbol(symbol byte) byte {
	// Letter set
	if (symbol >= 'a' && symbol <= 'z') || (symbol >= 'A' && symbol <= 'Z') {
		return 'L'
	}
	// Digit set
	if symbol >= '0' && symbol <= '9' {
		return 'D'
	}

	return symbol
}

// Scans a code file
func Scanner(file *os.File) Token {
	var lexemBuffer []byte
	buffer := make([]byte, 1)
	prevState := 0

	for {
		_, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF && len(lexemBuffer) == 0 {
			return EOF_TOKEN
		}

		currSymbol := RecognizeSymbol(buffer[0])
		currState, stateErr := GetNextState(prevState, currSymbol)

		// If the lexemBuffer has any symbol left we still need to consume this lexem
		if err == io.EOF && len(lexemBuffer) != 0 {
			tokenClass, err := GetTokenClass(currState)
			if err == ErrorNotFinalState {
				return ERROR_TOKEN
			}
			return NewToken(tokenClass, string(lexemBuffer), NULL)
		}

		if stateErr == ErrorTransitionDoesNotExist {
			tokenClass, err := GetTokenClass(currState)
			if err == ErrorNotFinalState {
				return ERROR_TOKEN
			}
			file.Seek(-1, 1)
			return NewToken(tokenClass, string(lexemBuffer), NULL)
		}

		lexemBuffer = append(lexemBuffer, buffer[0])
		prevState = currState
	}
}
