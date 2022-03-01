package parser

import (
	"mgol-go/src/lexer"
	"mgol-go/src/stack"
)

type Parser struct {
	stack *stack.Stack
	rules *RulesMap
	// actionTablePath string
}

func NewParser(stack *stack.Stack, rules *RulesMap, actionTablePath string) *Parser {
	return &Parser{
		stack: stack,
		rules: rules,
	}
}

func (p *Parser) Parse(token lexer.Token) {
	// algoritmo de fazer parser
	// state, _ := p.stack.Pop()
	// action := NewActionReader(p.actionTablePath)
	// act, opr := action.GetAction(state.(lexer.State), token)
	// switch act {

	// }
}
