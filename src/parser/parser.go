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

var errorsMessage = map[int]string{
	1: "token inesperado",
	2: "declaração de variáveis mal formada",
	3: "declaração de variáveis fora do escopo",
	4: "estrutura condicional mal formada",
	5: "estrutura de repetição mal formada",
	6: "tentativa de declaração inválida",
	7: "expressão inválida",
	8: "operação de entrada e saída inválida",
	9: "parênteses desbalanceados",
}

var parserErrorFlag = false

type Parser struct {
	scanner         *lexer.Scanner
	stack           *stack.Stack
	rules           *RulesMap
	semantic        *Semantic
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
		semantic:        NewSemantic(scanner.GetSymbolTable()),
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
	token, line, column := p.scanner.Scan()
	for isInTokensToIgnore(token) {
		token, line, column = p.scanner.Scan()
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
			p.semantic.semanticStack.Push(token)
			token, line, column = p.scanner.Scan()
			for isInTokensToIgnore(token) {
				token, line, column = p.scanner.Scan()
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
			p.semantic.ExecuteRule(rule, line, column)
		case ACCEPT:
			goto end_for
		case ERROR:
			errorMessage := getErrorMessage(opr)
			log.Printf("Erro: %v na linha %v, coluna %v", errorMessage, line, column)
			parserErrorFlag = true
			recoveryStatus := panicMode(p, token)

			if recoveryStatus == recoveryFail {
				goto end_for
			}
		}
	}
end_for:
	if semanticErrorFlag == false && parserErrorFlag == false {
		p.semantic.GenerateCode()
	}
	p.semantic.symbolTable.Print()
}

func getErrorMessage(id int) string {
	return errorsMessage[id]
}
