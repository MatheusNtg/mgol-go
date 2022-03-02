package lexer

import (
	"errors"
	"io"
	"log"
	errorhandling "mgol-go/src/error_handling"
	"os"
	"strings"
)

func letterGenerator() []Symbol {
	result := []Symbol{}

	for i := 'a'; i <= 'z'; i++ {
		result = append(result, Symbol(i))
	}

	for i := 'A'; i <= 'Z'; i++ {
		result = append(result, Symbol(i))
	}

	return result
}

func numGenerator() []Symbol {
	result := []Symbol{}

	for i := '0'; i <= '9'; i++ {
		result = append(result, Symbol(i))
	}

	return result
}

var (
	letters = letterGenerator()
	numbers = numGenerator()
)

func flatten(symbols [][]Symbol) []Symbol {
	result := []Symbol{}

	for _, arr := range symbols {
		result = append(result, arr...)
	}

	return result
}

var (
	alphabet = flatten([][]Symbol{
		letters,
		numbers,
		{
			'\n', '\t', ' ',
			'_', '+', '-', '*', '/',
			'>', '<', '=', '{', '}',
			'(', ')', ';', '"', '.',
			'E', 'e', ':', ',', '!',
			'?', '[', ']', '\\',
		},
	})
	states        = []State{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
	finalStates   = []State{1, 2, 4, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 20, 22, 25}
	transitionMap = map[State][]Transition{
		0: {
			{
				from: 0,
				to:   1,
				reading: flatten([][]Symbol{
					letters,
				}),
			},
			{
				from: 0,
				to:   2,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
			{
				from: 0,
				to:   8,
				reading: flatten([][]Symbol{
					{'<'},
				}),
			},
			{
				from: 0,
				to:   10,
				reading: flatten([][]Symbol{
					{'>'},
				}),
			},
			{
				from: 0,
				to:   12,
				reading: flatten([][]Symbol{
					{'='},
				}),
			},
			{
				from: 0,
				to:   14,
				reading: flatten([][]Symbol{
					{'+', '-', '*', '/'},
				}),
			},
			{
				from: 0,
				to:   15,
				reading: flatten([][]Symbol{
					{'('},
				}),
			},
			{
				from: 0,
				to:   16,
				reading: flatten([][]Symbol{
					{')'},
				}),
			},
			{
				from: 0,
				to:   17,
				reading: flatten([][]Symbol{
					{';'},
				}),
			},
			{
				from: 0,
				to:   19,
				reading: flatten([][]Symbol{
					{'{'},
				}),
			},
			{
				from: 0,
				to:   21,
				reading: flatten([][]Symbol{
					{'"'},
				}),
			},
		},

		1: {
			{
				from: 1,
				to:   1,
				reading: flatten([][]Symbol{
					letters,
					numbers,
					{'_'},
				}),
			},
		},

		2: {
			{
				from: 2,
				to:   2,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
			{
				from: 2,
				to:   3,
				reading: flatten([][]Symbol{
					{'.'},
				}),
			},
			{
				from: 2,
				to:   5,
				reading: flatten([][]Symbol{
					{'e', 'E'},
				}),
			},
		},

		3: {
			{
				from: 3,
				to:   4,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},

		4: {
			{
				from: 4,
				to:   4,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
			{
				from: 4,
				to:   23,
				reading: flatten([][]Symbol{
					{'e', 'E'},
				}),
			},
		},

		5: {
			{
				from: 5,
				to:   6,
				reading: flatten([][]Symbol{
					{'+', '-'},
				}),
			},
			{
				from: 5,
				to:   7,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},

		6: {
			{
				from: 6,
				to:   7,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},

		7: {
			{
				from: 7,
				to:   7,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},

		8: {
			{
				from: 8,
				to:   9,
				reading: flatten([][]Symbol{
					{'>', '='},
				}),
			},
			{
				from: 8,
				to:   13,
				reading: flatten([][]Symbol{
					{'-'},
				}),
			},
		},

		10: {
			{
				from: 10,
				to:   11,
				reading: flatten([][]Symbol{
					{'='},
				}),
			},
		},

		19: {
			{
				from: 19,
				to:   19,
				reading: flatten([][]Symbol{
					letters,
					numbers,
					{'\t', ' ', '_', '+', '-', '*', '/', '>', '<', '=', '{', '(', ')', ';', '"', '.', ':', ',', '!', '?', '[', ']', '\\'},
				}),
			},
			{
				from: 19,
				to:   20,
				reading: flatten([][]Symbol{
					{'}'},
				}),
			},
		},

		21: {
			{
				from: 21,
				to:   21,
				reading: flatten([][]Symbol{
					letters,
					numbers,
					{'\t', ' ', '_', '+', '-', '*', '/', '>', '<', '=', '{', '}', '(', ')', ';', '.', ':', ',', '!', '?', '[', ']', '\\'},
				}),
			},
			{
				from: 21,
				to:   22,
				reading: flatten([][]Symbol{
					{'"'},
				}),
			},
		},

		23: {
			{
				from: 23,
				to:   24,
				reading: flatten([][]Symbol{
					{'+', '-'},
				}),
			},
			{
				from: 23,
				to:   25,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},

		24: {
			{
				from: 24,
				to:   25,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},

		25: {
			{
				from: 25,
				to:   25,
				reading: flatten([][]Symbol{
					numbers,
				}),
			},
		},
	}
	stateToTokenClassMap = map[State]TokenClass{
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
		25: NUM,
	}
	numericTypes = map[State]DataType{
		2:  INTEGER,
		4:  REAL,
		7:  INTEGER,
		25: REAL,
	}
)

type Scanner struct {
	file                 *os.File
	lexemBuffer          []byte
	currentLineFile      int
	currentColumnFile    int
	dft                  Dft
	stateToTokenClassMap map[State]TokenClass
	symbolsToIgnore      []Symbol
	symbolTable          *SymbolTable
}

func NewScanner(file *os.File, symbolTable *SymbolTable) *Scanner {
	dft, err := NewDft(alphabet, states, 0, finalStates, transitionMap)
	if err != nil {
		log.Fatal("Failed to create DFT:", err)
	}

	return &Scanner{
		file:                 file,
		lexemBuffer:          []byte{},
		currentLineFile:      1,
		currentColumnFile:    0,
		dft:                  *dft,
		stateToTokenClassMap: stateToTokenClassMap,
		symbolsToIgnore:      []Symbol{'\n', ' ', '\t'},
		symbolTable:          symbolTable,
	}
}

func (s *Scanner) getTokenClass(state State) TokenClass {
	return s.stateToTokenClassMap[state]
}

// Recognizes if a number is integer or real and
// changes the dataType accordingly
func (s *Scanner) updateDataType(token *Token) {
	switch token.class {
	case NUM:
		token.dataType = numericTypes[s.dft.currentState]
	case LITERAL_CONST:
		token.dataType = LITERAL
	default:
		token.dataType = NULL
	}
}

func (s *Scanner) clearLexemBuffer() {
	s.lexemBuffer = []byte{}
}

// reset clears the lexem buffer
// and resets the head of the dft
func (s *Scanner) reset() {
	s.clearLexemBuffer()
	s.dft.Reset()
}

// resetAndRewind does the same as
// reset but rewind the head of the
// file as well
func (s *Scanner) resetAndRewind() {
	s.reset()
	s.file.Seek(-1, os.SEEK_CUR)
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
		currSymbol := Symbol(currChar)

		s.currentColumnFile += n

		if err == io.EOF && len(s.lexemBuffer) == 0 {
			return EOF_TOKEN
		}

		if err == io.EOF && len(s.lexemBuffer) != 0 {
			if ContainsByte(s.lexemBuffer, '{') && !ContainsByte(s.lexemBuffer, '}') {
				errorhandling.NewLexicalError(s.currentLineFile, s.currentColumnFile, string(s.lexemBuffer))
				s.reset()
				return ERROR_TOKEN
			}

			numberOfQuotation := strings.Count(string(s.lexemBuffer), "\"")
			if numberOfQuotation == 1 {
				errorhandling.NewLexicalError(s.currentLineFile, s.currentColumnFile, string(s.lexemBuffer))
				s.reset()
				return ERROR_TOKEN
			}

			tokenClass := s.getTokenClass(s.dft.GetCurrentState())
			if tokenClass == COMMENT {
				s.reset()
				return COMMENT_TOKEN
			}
			lexem := s.lexemBuffer
			token := NewToken(tokenClass, string(lexem), NULL)
			s.updateDataType(&token)

			s.reset()

			if token.class == IDENTIFIER {
				return s.symbolTable.Insert(token.lexeme, token)
			}
			return token
		}

		if !ContainsSymbol(alphabet, currSymbol) || !ContainsByte(s.lexemBuffer, '{') && currChar == '}' {
			errorhandling.NewLexicalError(s.currentLineFile, s.currentColumnFile, string(s.lexemBuffer)+string(currChar))
			s.reset()
			return ERROR_TOKEN
		}

		if currChar == '\n' {
			s.currentLineFile += 1
			s.currentColumnFile = 0
		}

		_, err = s.dft.Next(currSymbol)

		if errors.Is(err, ErrorTransitionDoesNotExist) && s.dft.IsFinalState() {
			tokenClass := s.getTokenClass(s.dft.GetCurrentState())
			if tokenClass == COMMENT {
				s.resetAndRewind()
				return COMMENT_TOKEN
			}
			lexem := s.lexemBuffer
			token := NewToken(tokenClass, string(lexem), NULL)
			s.updateDataType(&token)

			s.resetAndRewind()

			s.currentColumnFile -= n
			if currChar == '\n' {
				s.currentLineFile -= 1
			}

			if token.class == IDENTIFIER {
				return s.symbolTable.Insert(token.lexeme, token)
			}
			return token
		}

		if errors.Is(err, ErrorTransitionDoesNotExist) && !s.dft.IsFinalState() {
			if currChar == ' ' || currChar == '\n' || currChar == '\t' {
				continue
			}

			if len(string(s.lexemBuffer)) == 0 {
				errorhandling.NewLexicalError(s.currentLineFile, s.currentColumnFile, string(currChar))
			} else {
				errorhandling.NewLexicalError(s.currentLineFile, s.currentColumnFile, string(s.lexemBuffer))
			}

			s.clearLexemBuffer()
			if s.dft.currentState != s.dft.initialState {
				s.file.Seek(-1, os.SEEK_CUR)
			}
			s.dft.Reset()

			return ERROR_TOKEN
		}

		if !ContainsSymbol(s.symbolsToIgnore, currSymbol) {
			s.lexemBuffer = append(s.lexemBuffer, currChar)
		} else if ContainsByte(s.lexemBuffer, '"') || ContainsByte(s.lexemBuffer, '{') {
			s.lexemBuffer = append(s.lexemBuffer, currChar)
		}
	}
}
