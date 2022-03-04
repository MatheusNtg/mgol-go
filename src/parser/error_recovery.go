package parser

import (
	"log"
	"mgol-go/src/lexer"
)

type RecoveryStatus bool

const (
	recoverySucess RecoveryStatus = true
	recoveryFail   RecoveryStatus = false
)

func panicMode(parser *Parser, firstToken lexer.Token, line, column int) RecoveryStatus {

	stack_copy := parser.stack.Clone()
	actionReader := NewActionReader(parser.actionTablePath)
	token := firstToken
	log.Printf("Token \"%v\" inesperado na linha %v coluna %v\n", token.GetLexem(), line, column)
	parser.stack.Pop()

	for {
		for parser.stack.GetLength() > 0 {
			topStack, err := parser.stack.Get()
			if err != nil {
				panic(err)
			}

			state := lexer.State(topStack.(int))
			action, _ := actionReader.GetAction(state, token)

			if action != NONE {
				return recoverySucess
			}
			parser.stack.Pop()
		}

		parser.stack = stack_copy
		token, line, column = parser.scanner.Scan()
		log.Printf("Token \"%v\" inesperado na linha %v coluna %v\n", token.GetLexem(), line, column)
		if token == lexer.EOF_TOKEN {
			return recoveryFail
		}
	}
}
