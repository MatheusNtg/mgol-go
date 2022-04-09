package parser

import (
	"fmt"
	"io/ioutil"
	"mgol-go/src/lexer"
	"mgol-go/src/stack"
)

type TemporalType int

const (
	TemporalBool TemporalType = iota
	TemporalInt
)

const maxCapacityStack = 10000

type CodeBuffer struct {
	temporals []TemporalType
	code      string
}

func NewCodeBuffer() *CodeBuffer {
	return &CodeBuffer{}
}

var rulesMap = map[int]func(s *Semantic, rule Rule){
	6: func(s *Semantic, rule Rule) {
		s.AddToCodeBuffer(";\n")
	},

	7: func(s *Semantic, rule Rule) {
		identifierToken, _ := s.semanticStack.Pop()
		identifierTokenConverted := identifierToken.(lexer.Token)

		typeToken, _ := s.semanticStack.Pop()
		typeTokenConverted := typeToken.(lexer.Token)

		identifierTokenConverted.SetType(typeTokenConverted.GetType())
		s.symbolTable.Update(identifierTokenConverted.GetLexem(), identifierTokenConverted)

		s.AddToCodeBuffer(identifierTokenConverted.GetLexem())
	},

	8: func(s *Semantic, rule Rule) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.INTEGER)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("int ")
	},

	9: func(s *Semantic, rule Rule) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.REAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("float ")
	},

	10: func(s *Semantic, rule Rule) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.LITERAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("literal ")
	},

	// ES -> leia id pt_v
	12: func(s *Semantic, rule Rule) {
		s.semanticStack.Pop() // Remove our pt_v
		idToken, _ := s.semanticStack.Pop()
		idTokenConverted := idToken.(lexer.Token)
		if idTokenConverted.GetType() == lexer.NULL {
			//TODO Improve our error handlings
			fmt.Printf("Erro de variável não declarada")
			return
		}
		switch idTokenConverted.GetType() {
		case lexer.INTEGER:
			s.AddToCodeBuffer(fmt.Sprintf("scanf(\"%%d\", %s);\n", idTokenConverted.GetLexem()))
		case lexer.LITERAL:
			s.AddToCodeBuffer(fmt.Sprintf("scanf(\"%%s\", &%s);\n", idTokenConverted.GetLexem()))
		case lexer.REAL:
			s.AddToCodeBuffer(fmt.Sprintf("scanf(\"%%lf\", &%s);\n", idTokenConverted.GetLexem()))
		}
	},

	13: func(s *Semantic, rule Rule) {
		s.semanticStack.Pop() // Remove our pt_v
		argToken, _ := s.semanticStack.Pop()
		argTokenConverted := argToken.(lexer.Token)
		switch argTokenConverted.GetType() {
		case lexer.INTEGER:
			s.AddToCodeBuffer(fmt.Sprintf("printf(\"%%d\", %s);\n", argTokenConverted.GetLexem()))
		case lexer.LITERAL:
			s.AddToCodeBuffer(fmt.Sprintf("printf(\"%%s\", &%s);\n", argTokenConverted.GetLexem()))
		case lexer.REAL:
			s.AddToCodeBuffer(fmt.Sprintf("printf(\"%%lf\", &%s);\n", argTokenConverted.GetLexem()))
		}
	},

	14: func(s *Semantic, rule Rule) {
		literalToken, _ := s.symbolTable.GetToken("literal")
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), literalToken.GetLexem(), literalToken.GetType())
		s.semanticStack.Push(newToken)
	},

	15: func(s *Semantic, rule Rule) {
		numToken, _ := s.semanticStack.Pop()
		numTokenConverted := numToken.(lexer.Token)
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), numTokenConverted.GetLexem(), numTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	21: func(s *Semantic, rule Rule) {
		idToken, _ := s.semanticStack.Pop()
		idTokenConverted := idToken.(lexer.Token)
		if idTokenConverted.GetType() == lexer.NULL {
			//TODO Improve our error handlings
			fmt.Printf("Erro de variável não declarada")
			return
		}
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), idTokenConverted.GetLexem(), idTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	22: func(s *Semantic, rule Rule) {
		numToken, _ := s.semanticStack.Pop()
		numTokenConverted := numToken.(lexer.Token)

		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), numTokenConverted.GetLexem(), numTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// COND -> CAB CP
	24: func(s *Semantic, rule Rule) {
		s.AddToCodeBuffer("}")
	},

	// CAB -> se ab_p EXP_R fc_p entao
	25: func(s *Semantic, rule Rule) {
		s.AddToCodeBuffer("if")
	},

	// EXP_R -> OPRD opr OPRD
	26: func(s *Semantic, rule Rule) {
		rawOprd2, _ := s.semanticStack.Pop()
		oprd2 := rawOprd2.(lexer.Token)

		rawOpr, _ := s.semanticStack.Pop()
		opr := rawOpr.(lexer.Token)

		rawOprd1, _ := s.semanticStack.Pop()
		oprd1 := rawOprd1.(lexer.Token)

		if oprd1.GetType() != oprd2.GetType() {
			fmt.Print("Erro: Operandos com tipos incompativeis")
		}

		temporalId := s.NewTemporal(TemporalBool)

		exp_rToken := lexer.NewToken(lexer.TokenClass(rule.Left), temporalId, lexer.NULL)
		s.semanticStack.Push(exp_rToken)
		s.AddToCodeBuffer(fmt.Sprintf("%s = %s %s %s", temporalId, oprd1.GetLexem(), opr.GetType(), oprd2.GetLexem()))
	},
}

type Semantic struct {
	semanticStack *stack.Stack
	codeBuffer    *CodeBuffer
	ruleMap       map[int]func(s *Semantic, rule Rule)
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

func (s *Semantic) ExecuteRule(rule Rule) {
	_, found := s.ruleMap[rule.Number+1]
	if !found {
		return
	}
	s.ruleMap[rule.Number+1](s, rule)
}

func (s *Semantic) AddToCodeBuffer(code string) {
	s.codeBuffer.code += code
}

// NewTemporal adds a new temporal variable of TemporalType
func (s *Semantic) NewTemporal(temporalType TemporalType) string {
	s.codeBuffer.temporals = append(s.codeBuffer.temporals, temporalType)
	return fmt.Sprintf("T%d", (s.codeBuffer.temporals))
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
