package parser

import (
	"mgol-go/src/lexer"
)

type RecoveryStatus bool

const (
	recoverySucess RecoveryStatus = true
	recoveryFail   RecoveryStatus = false
)

func panicMode(parser *Parser, firstToken lexer.Token) RecoveryStatus {

	stack_copy := parser.stack.Clone()
	actionReader := NewActionReader(parser.actionTablePath)
	token := firstToken
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
		token = parser.scanner.Scan()
		if token == lexer.EOF_TOKEN {
			return recoveryFail
		}
	}
}
