package parser

import (
	"fmt"
	"io/ioutil"
	"mgol-go/src/lexer"
	"mgol-go/src/stack"
)

const maxCapacityStack = 10000

type CodeBuffer struct {
	// temporals []string
	code string
}

func NewCodeBuffer() *CodeBuffer {
	return &CodeBuffer{}
}

var rulesMap = map[int]func(s *Semantic, rule Rule, token lexer.Token){
	5: func(s *Semantic, rule Rule, token lexer.Token) {
		topStackToken, err := s.semanticStack.Pop()
		if err != nil {
			panic(err)
		}
		convertedToken := topStackToken.(lexer.Token)
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", convertedToken.GetType())
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer(";\n")
		s.semanticStack.Push(newToken)
	},

	6: func(s *Semantic, rule Rule, token lexer.Token) {
		topStackToken, err := s.semanticStack.Pop()
		if err != nil {
			panic(err)
		}
		convertedToken := topStackToken.(lexer.Token)
		newToken := lexer.NewToken(lexer.TokenClass(rule.Right[0]), convertedToken.GetLexem(), convertedToken.GetType())
		s.symbolTable.Update(newToken.GetLexem(), newToken)
		s.AddToCodeBuffer(newToken.GetLexem())
		s.semanticStack.Push(newToken)
	},

	7: func(s *Semantic, rule Rule, token lexer.Token) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), token.GetLexem(), lexer.INTEGER)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("int ")
	},

	8: func(s *Semantic, rule Rule, token lexer.Token) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), token.GetLexem(), lexer.REAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("float ")
	},

	9: func(s *Semantic, rule Rule, token lexer.Token) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), token.GetLexem(), lexer.LITERAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("literal ")
	},
}

type Semantic struct {
	semanticStack *stack.Stack
	codeBuffer    *CodeBuffer
	ruleMap       map[int]func(s *Semantic, rule Rule, token lexer.Token)
	symbolTable   *lexer.SymbolTable
}

func NewSemantic(symbolTable *lexer.SymbolTable) *Semantic {
	return &Semantic{
		semanticStack: stack.NewStack(maxCapacityStack),
		codeBuffer:    NewCodeBuffer(),
		ruleMap:       rulesMap,
		symbolTable:   symbolTable,
	}
}

func (s *Semantic) ExecuteRule(rule Rule, token lexer.Token) {
	_, found := s.ruleMap[rule.Number]
	if !found {
		return
	}
	s.ruleMap[rule.Number](s, rule, token)
}

func (s *Semantic) AddToCodeBuffer(code string) {
	s.codeBuffer.code += code
}

func (s *Semantic) GenerateCode() {
	currentCode := `
#include<stdio.h>
typedef char literal[256];
void main() {
`

	currentCode = fmt.Sprintf("%s%s", currentCode, s.codeBuffer.code)

	currentCode = fmt.Sprintf("%s%s", currentCode, "\n}")

	ioutil.WriteFile("output.c", []byte(currentCode), 0755)
}
