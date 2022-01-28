package errorhandling

import (
	"log"
	"strings"
)

type LexicalErrorType int

const (
	InvalidLiteral LexicalErrorType = iota
	InvalidNumber
	InvalidComment
	InvalidWord
)

func isInvalidNumber(lexem string) bool {
	containsQuotation := strings.Contains(lexem, "\"")
	containsBrackets := strings.Contains(lexem, "{")
	containsNumber := strings.ContainsAny(lexem, "0123456789")

	if !containsBrackets && !containsQuotation && containsNumber {
		return true
	}
	return false
}

func isInvalidLiteral(lexem string) bool {
	numberOfQuotations := strings.Count(lexem, "\"")
	indexOfQuotation := strings.Index(lexem, "\"")
	if numberOfQuotations == 1 && indexOfQuotation == 0 {
		return true
	}
	return false
}

func isInvalidComment(lexem string) bool {
	numberOfBrackets := strings.Count(lexem, "{")
	indexOfBracket := strings.Index(lexem, "{")
	if numberOfBrackets == 1 && indexOfBracket == 0 {
		return true
	}
	return false
}

func getErrorType(lexem string) LexicalErrorType {
	if isInvalidNumber(lexem) {
		return InvalidNumber
	}

	if isInvalidLiteral(lexem) {
		return InvalidLiteral
	}

	if isInvalidComment(lexem) {
		return InvalidComment
	}

	return InvalidWord
}

func NewLexicalError(line, column int, lexem string) {
	errorType := getErrorType(lexem)

	switch errorType {
	case InvalidLiteral:
		log.Printf("erro na linha %d coluna %d, literal %s inválido", line, column, lexem)
	case InvalidNumber:
		log.Printf("erro na linha %d coluna %d, número %s inválido", line, column, lexem)
	case InvalidComment:
		log.Printf("erro na linha %d coluna %d, comentário %s inválido", line, column, lexem)
	case InvalidWord:
		log.Printf("erro na linha %d coluna %d, palavra %s inexistente na linguagem", line, column, lexem)
	}
}
