package parser

import (
	"log"
	"mgol-go/src/lexer"
	"mgol-go/src/stack"
)

type Parser struct {
	scanner         *lexer.Scanner
	stack           *stack.Stack
	rules           *RulesMap
	actionTablePath string
	gotoTablePath   string
}

func NewParser(scanner *lexer.Scanner, stack *stack.Stack, rules *RulesMap, actionTablePath, gotoTablePath string) *Parser {
	return &Parser{
		scanner:         scanner,
		stack:           stack,
		rules:           rules,
		actionTablePath: actionTablePath,
		gotoTablePath:   gotoTablePath,
	}
}

func (p *Parser) Parse() {
	token := p.scanner.Scan()
	p.stack.Push(0)

	actionReader := NewActionReader(p.actionTablePath)
	gotoReader := NewGotoReader(p.gotoTablePath)

	for {
		topStack, err := p.stack.Get()
		if err != nil {
			panic(err)
		}

		state := lexer.State(topStack.(int))
		action, opr := actionReader.GetAction(state, token)
		switch action {
		case SHIFT:
			p.stack.Push(opr)
			token = p.scanner.Scan()
		case REDUCE:
			rule := p.rules.GetRule(opr)
			for range rule.Right {
				p.stack.Pop()
			}
			p.stack.Push(opr)
			currentTopElement, err := p.stack.Get()
			state = lexer.State(currentTopElement.(int))
			if err != nil {
				panic(err)
			}
			gotoOpr := gotoReader.GetGoto(state, rule.Left)
			p.stack.Push(gotoOpr)
			log.Print(rule.Left, "->", rule.Right)
		case ACCEPT:
			goto for_end
		case NONE:
			log.Println("Deu pau")
		}
	}
for_end:
}
