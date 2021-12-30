package lexer

import (
	"errors"
	"io"
	"log"
	"os"
)

var (
	alphabet      = []byte{'\n', '\t', ' ', 'L', 'D', '_'}
	states        = []int{0, 1}
	finalStates   = []int{1}
	transitionMap = map[int]map[byte]int{
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
	stateToTokenClassMap = map[int]TokenClass{
		1: IDENTIFIER,
	}
)

func containsByte(slice []byte, element byte) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

type Scanner struct {
	file                 *os.File
	lexemBuffer          []byte
	currentLineFile      int
	currentColumnFile    int
	dft                  Dft
	stateToTokenClassMap map[int]TokenClass
	charsToIgnore        []byte
}

func NewScanner(file *os.File) *Scanner {
	dft, err := NewDft(alphabet, states, 0, finalStates, transitionMap)
	if err != nil {
		log.Fatal("Failed to create DFT:", err)
	}

	return &Scanner{
		file:                 file,
		lexemBuffer:          []byte{},
		currentLineFile:      1,
		currentColumnFile:    1,
		dft:                  *dft,
		stateToTokenClassMap: stateToTokenClassMap,
		charsToIgnore:        []byte{'\n', ' ', '\t'},
	}
}

func (s *Scanner) getTokenClass(state int) TokenClass {
	return s.stateToTokenClassMap[state]
}

// TODO: Implement this method to recognize the datatypes
func (s *Scanner) getDataType(lexem string) DataType {
	return NULL
}

func (s *Scanner) clearLexemBuffer() {
	s.lexemBuffer = []byte{}
}

func (s *Scanner) recognizeSymbol(symbol byte) byte {
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

// Scan reads the Scanner file until finds an Token or an error.
// If it finds an Token it returns the reconized token, otherwhise
// just returns an error Token and shows to the user the error
// message related
func (s *Scanner) Scan() Token {
	readBuffer := make([]byte, 1)

	for {
		n, err := s.file.Read(readBuffer)
		currChar := readBuffer[0]
		currSymbol := s.recognizeSymbol(currChar)

		s.currentColumnFile += n
		if err == io.EOF && len(s.lexemBuffer) == 0 {
			return EOF_TOKEN
		}

		if err == io.EOF && len(s.lexemBuffer) != 0 {
			tokenClass := s.getTokenClass(s.dft.GetCurrentState())
			lexem := s.lexemBuffer
			dataType := s.getDataType(string(lexem))

			s.clearLexemBuffer()

			return NewToken(tokenClass, string(lexem), dataType)
		}

		if currChar == '\n' {
			s.currentLineFile += 1
			s.currentColumnFile = 1
		}

		if containsByte(s.charsToIgnore, currChar) {
			s.dft.Next(currSymbol)
			continue
		}

		_, err = s.dft.Next(currSymbol)
		//TODO: Tratar o erro aqui
		if errors.Is(err, ErrorTransitionDoesNotExist) && s.dft.IsFinalState() {
			tokenClass := s.getTokenClass(s.dft.GetCurrentState())
			lexem := s.lexemBuffer
			dataType := s.getDataType(string(lexem))

			s.clearLexemBuffer()
			s.file.Seek(-1, os.SEEK_CUR)

			return NewToken(tokenClass, string(lexem), dataType)
		}

		s.lexemBuffer = append(s.lexemBuffer, currChar)
	}
}
