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
		s.AddToCodeBuffer(";\n")
	},

	6: func(s *Semantic, rule Rule, token lexer.Token) {
		identifierToken, _ := s.semanticStack.Pop()
		identifierTokenConverted := identifierToken.(lexer.Token)

		typeToken, _ := s.semanticStack.Pop()
		typeTokenConverted := typeToken.(lexer.Token)

		identifierTokenConverted.SetType(typeTokenConverted.GetType())
		s.symbolTable.Update(identifierTokenConverted.GetLexem(), identifierTokenConverted)

		s.AddToCodeBuffer(identifierTokenConverted.GetLexem())
	},

	7: func(s *Semantic, rule Rule, token lexer.Token) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.INTEGER)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("int ")
		s.semanticStack.Push(newToken)
	},

	8: func(s *Semantic, rule Rule, token lexer.Token) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.REAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("float ")
		s.semanticStack.Push(newToken)
	},

	9: func(s *Semantic, rule Rule, token lexer.Token) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.LITERAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("literal ")
		s.semanticStack.Push(newToken)
	},

	11: func(s *Semantic, rule Rule, token lexer.Token) {
		fmt.Print(rule, token)
	},

	// 12: func(s *Semantic, rule Rule, token lexer.Token) {

	// },
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
