package parser

import "mgol-go/src/stack"

type Parser struct {
	stack *stack.Stack
	rules *RulesMap
}

func NewParser(stack *stack.Stack, rules *RulesMap) *Parser {
	return &Parser{
		stack: stack,
		rules: rules,
	}
}

// func Parse()
