package lexer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	alphabet = []byte{
		'\n', '\t', ' ', 'L', 'D',
		'_', '+', '-', '*', '/',
		'>', '<', '=', '{', '}',
		'(', ')', ';', '"', '.',
		'E', 'e', ':', ',', '!',
		'?', '[', ']', '\\',
	}
	states        = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22}
	finalStates   = []int{1, 2, 4, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 20, 22}
	transitionMap = map[int]map[byte]int{
		0: {
			'\n': 0,
			'\t': 0,
			' ':  0,
			'L':  1, // L is the set that represents all the letters of the alphabet, low and capital case
			'D':  2,
			'<':  8,
			'>':  10,
			'=':  12, // Relational operator
			'+':  14, // Arithmitec Operator
			'-':  14,
			'*':  14,
			'/':  14,
			'(':  15, // Open brackets
			')':  16, // Closed brackets
			';':  17, // Semi columns
			'{':  19,
			'"':  21,
		},
		// Identifier
		1: {
			'L': 1,
			'D': 1, // D is the set that represents all numbers from 0 to 9
			'_': 1,
		},
		// Numeric Constants
		2: {
			'D': 2,
			'.': 3,
			'E': 5,
			'e': 5,
		},
		3: {
			'D': 4,
		},
		4: {
			'D': 4,
			'E': 5,
			'e': 5,
		},
		5: {
			'+': 6,
			'-': 6,
			'D': 7,
		},
		6: {
			'D': 7,
		},
		7: {
			'D': 7,
		},
		// Relational operator
		8: {
			'>': 9,
			'=': 9,
			'-': 13, // Attribute
		},
		10: {
			'=': 11,
		},
		19: {
			'\t': 19,
			'\n': 19,
			' ':  19,
			'L':  19,
			'D':  19,
			'_':  19,
			'+':  19,
			'-':  19,
			'*':  19,
			'/':  19,
			'>':  19,
			'<':  19,
			'=':  19,
			'{':  19,
			'(':  19,
			')':  19,
			';':  19,
			'"':  19,
			'.':  19,
			'E':  19,
			'e':  19,
			':':  19,
			',':  19,
			'!':  19,
			'?':  19,
			'[':  19,
			']':  19,
			'\\': 19,
			'}':  20,
		},
		21: {
			'\t': 21,
			'\n': 21,
			' ':  21,
			'L':  21,
			'D':  21,
			'_':  21,
			'+':  21,
			'-':  21,
			'*':  21,
			'/':  21,
			'>':  21,
			'<':  21,
			'=':  21,
			'{':  21,
			'(':  21,
			')':  21,
			';':  21,
			'.':  21,
			'E':  21,
			'e':  21,
			'}':  21,
			':':  21,
			',':  21,
			'!':  21,
			'?':  21,
			'[':  21,
			']':  21,
			'\\': 21,
			'"':  22,
		},
	}
	stateToTokenClassMap = map[int]TokenClass{
		1:  IDENTIFIER,
		2:  NUM,
		4:  NUM,
		7:  NUM,
		8:  REL_OP,
		9:  REL_OP,
		10: REL_OP,
		11: REL_OP,
		12: REL_OP,
		13: ATTR,
		14: ARIT_OP,
		15: OPEN_PAR,
		16: CLOSE_PAR,
		17: SEMICOLON,
		20: COMMENT,
		22: LITERAL_CONST,
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

// Scan reads the Scanner file until finds a Token or an error.
// If it finds a Token it returns the reconized token, otherwhise
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
			s.dft.Reset()
			token := NewToken(tokenClass, string(lexem), dataType)
			InsertSymbolTable(string(lexem), token)

			return token
		}

		if !containsByte(alphabet, currSymbol) {
			log.Printf("Erro na linha %d coluna %d, palavra %s não existe na linguagem\n", s.currentLineFile, s.currentColumnFile, fmt.Sprintf("%s%s", string(s.lexemBuffer), string(currChar)))
			s.clearLexemBuffer()
			s.dft.Reset()
			return ERROR_TOKEN
		}

		if currChar == '\n' {
			s.currentLineFile += 1
			s.currentColumnFile = 1
		}

		if containsByte(s.charsToIgnore, currChar) {
			s.dft.Next(currSymbol)
			if s.dft.IsFinalState() {
				tokenClass := s.getTokenClass(s.dft.GetCurrentState())
				lexem := s.lexemBuffer
				dataType := s.getDataType(string(lexem))

				s.clearLexemBuffer()
				s.dft.Reset()
				s.file.Seek(-1, os.SEEK_CUR)

				token := NewToken(tokenClass, string(lexem), dataType)
				InsertSymbolTable(string(lexem), token)

				return token
			}
		}

		_, err = s.dft.Next(currSymbol)

		if errors.Is(err, ErrorTransitionDoesNotExist) && s.dft.IsFinalState() {
			tokenClass := s.getTokenClass(s.dft.GetCurrentState())
			lexem := s.lexemBuffer
			dataType := s.getDataType(string(lexem))

			s.clearLexemBuffer()
			s.dft.Reset()
			s.file.Seek(-1, os.SEEK_CUR)

			token := NewToken(tokenClass, string(lexem), dataType)
			InsertSymbolTable(string(lexem), token)

			return token
		}

		if !containsByte(s.charsToIgnore, currChar) {
			s.lexemBuffer = append(s.lexemBuffer, currChar)
		} else if containsByte(s.lexemBuffer, '"') {
			s.lexemBuffer = append(s.lexemBuffer, currChar)
		}
	}
}
