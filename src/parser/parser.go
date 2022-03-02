package parser

import (
	"log"
	"mgol-go/src/lexer"
	"mgol-go/src/stack"
)

var (
	tokensToIgnore = []lexer.Token{
		lexer.ERROR_TOKEN,
		lexer.COMMENT_TOKEN,
	}
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

// isInTokensToIgnore return whether a token
// t is in the list of tokens to ignore or not
func isInTokensToIgnore(t lexer.Token) bool {
	for _, token := range tokensToIgnore {
		if t == token {
			return true
		}
	}
	return false
}

func (p *Parser) Parse() {
	token := p.scanner.Scan()
	for isInTokensToIgnore(token) {
		token = p.scanner.Scan()
	}
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
			for isInTokensToIgnore(token) {
				token = p.scanner.Scan()
			}
		case REDUCE:
			rule := p.rules.GetRule(opr)
			for range rule.Right {
				p.stack.Pop()
			}
			currentTopElement, err := p.stack.Get()
			state = lexer.State(currentTopElement.(int))
			if err != nil {
				panic(err)
			}
			gotoOpr := gotoReader.GetGoto(state, rule.Left)
			p.stack.Push(gotoOpr)
			log.Print(rule.Left, "->", rule.Right)
		case ACCEPT:
			goto end_for
		case NONE:
			recoveryStatus := panicMode(p, token)

			if recoveryStatus == recoveryFail {
				goto end_for
			}
		}
	}
end_for:
}
