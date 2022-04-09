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
	TemporalFloat
)

const maxCapacityStack = 10000

var semanticErrorFlag = false

type CodeBuffer struct {
	temporals []TemporalType
	code      string
}

func NewCodeBuffer() *CodeBuffer {
	return &CodeBuffer{}
}

func (c *CodeBuffer) PrintTemporals() string {
	temporalCode := "/*----Variaveis temporarias----*/\n"
	for idx, temporal := range c.temporals {
		chunk := ""
		switch temporal {
		case TemporalBool:
			chunk = fmt.Sprintf("bool T%d;\n", idx)
		case TemporalInt:
			chunk = fmt.Sprintf("int T%d;\n", idx)
		case TemporalFloat:
			chunk = fmt.Sprintf("float T%d;\n", idx)
		}
		temporalCode += chunk
	}
	temporalCode += "/*------------------------------*/\n"
	return temporalCode
}

var rulesMap = map[int]func(s *Semantic, rule Rule, line int, column int){
	// D -> TIPO L pt_v
	6: func(s *Semantic, rule Rule, line int, column int) {
		s.AddToCodeBuffer(";\n")
	},

	// L -> id
	7: func(s *Semantic, rule Rule, line int, column int) {
		identifierToken, _ := s.semanticStack.Pop()
		identifierTokenConverted := identifierToken.(lexer.Token)

		typeToken, _ := s.semanticStack.Pop()
		typeTokenConverted := typeToken.(lexer.Token)

		identifierTokenConverted.SetType(typeTokenConverted.GetType())
		s.symbolTable.Update(identifierTokenConverted.GetLexem(), identifierTokenConverted)

		s.AddToCodeBuffer(identifierTokenConverted.GetLexem())
	},

	// TIPO -> inteiro
	8: func(s *Semantic, rule Rule, line int, column int) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.INTEGER)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("int ")
	},

	// TIPO -> real
	9: func(s *Semantic, rule Rule, line int, column int) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.REAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("float ")
	},

	// TIPO -> literal
	10: func(s *Semantic, rule Rule, line int, column int) {
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), "", lexer.LITERAL)
		s.semanticStack.Push(newToken)
		s.AddToCodeBuffer("literal ")
	},

	// ES -> leia id pt_v
	12: func(s *Semantic, rule Rule, line int, column int) {
		s.semanticStack.Pop() // Remove our pt_v
		idToken, _ := s.semanticStack.Pop()
		idTokenConverted := idToken.(lexer.Token)
		if idTokenConverted.GetType() == lexer.NULL {
			fmt.Printf("Erro: variável %s não declarada na linha %d, coluna %d\n", idTokenConverted.GetLexem(), line, column)
			semanticErrorFlag = true
			return
		}
		switch idTokenConverted.GetType() {
		case lexer.INTEGER:
			s.AddToCodeBuffer(fmt.Sprintf("scanf(\"%%d\", &%s);\n", idTokenConverted.GetLexem()))
		case lexer.LITERAL:
			s.AddToCodeBuffer(fmt.Sprintf("scanf(\"%%s\", %s);\n", idTokenConverted.GetLexem()))
		case lexer.REAL:
			s.AddToCodeBuffer(fmt.Sprintf("scanf(\"%%lf\", &%s);\n", idTokenConverted.GetLexem()))
		}
	},

	// ES -> escreva ARG pt_v
	13: func(s *Semantic, rule Rule, line int, column int) {
		s.semanticStack.Pop() // Remove our pt_v
		argToken, _ := s.semanticStack.Pop()
		argTokenConverted := argToken.(lexer.Token)
		switch argTokenConverted.GetType() {
		case lexer.INTEGER:
			s.AddToCodeBuffer(fmt.Sprintf("printf(\"%%d\", %s);\n", argTokenConverted.GetLexem()))
		case lexer.LITERAL:
			s.AddToCodeBuffer(fmt.Sprintf("printf(\"%%s\", %s);\n", argTokenConverted.GetLexem()))
		case lexer.REAL:
			s.AddToCodeBuffer(fmt.Sprintf("printf(\"%%lf\", %s);\n", argTokenConverted.GetLexem()))
		}
	},

	// ARG -> lit
	14: func(s *Semantic, rule Rule, line int, column int) {
		literalToken, _ := s.semanticStack.Pop()
		literalTokenConverted := literalToken.(lexer.Token)
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), literalTokenConverted.GetLexem(), literalTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// ARG -> num
	15: func(s *Semantic, rule Rule, line int, column int) {
		numToken, _ := s.semanticStack.Pop()
		numTokenConverted := numToken.(lexer.Token)
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), numTokenConverted.GetLexem(), numTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// ARG -> id
	16: func(s *Semantic, rule Rule, line int, column int) {
		idToken, _ := s.semanticStack.Pop()
		idTokenConverted := idToken.(lexer.Token)
		if idTokenConverted.GetType() == lexer.NULL {
			fmt.Printf("Erro: variável %s não declarada na linha %d, coluna %d\n", idTokenConverted.GetLexem(), line, column)
			semanticErrorFlag = true
			return
		}

		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), idTokenConverted.GetLexem(), idTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// CMD -> id rcb LD pt_v
	18: func(s *Semantic, rule Rule, line int, column int) {
		s.semanticStack.Pop() // remove our pt_v
		rawLD, _ := s.semanticStack.Pop()
		LD := rawLD.(lexer.Token)

		s.semanticStack.Pop() // remove our rcb

		rawId, _ := s.semanticStack.Pop()
		id := rawId.(lexer.Token)

		if id.GetType() == lexer.NULL {
			fmt.Printf("Erro: variável %s não declarada na linha %d, coluna %d\n", id.GetLexem(), line, column)
			semanticErrorFlag = true
			return
		}

		if id.GetType() != LD.GetType() {
			fmt.Printf("Erro: Tipos diferentes para a atribuição na linha %d, coluna %d. '%s' é do tipo '%s', enquanto que '%s' é do tipo '%s'\n", line, column, id.GetLexem(), id.GetType(), LD.GetLexem(), LD.GetType())
			semanticErrorFlag = true
			return
		}

		s.AddToCodeBuffer(fmt.Sprintf("%s = %s;\n", id.GetLexem(), LD.GetLexem()))
	},

	// LD -> OPRD opm OPRD
	19: func(s *Semantic, rule Rule, line int, column int) {
		rawOprd2, _ := s.semanticStack.Pop()
		oprd2 := rawOprd2.(lexer.Token)

		rawOpm, _ := s.semanticStack.Pop()
		opm := rawOpm.(lexer.Token)

		rawOprd1, _ := s.semanticStack.Pop()
		oprd1 := rawOprd1.(lexer.Token)

		if oprd1.GetType() != oprd2.GetType() && oprd1.GetType() != lexer.LITERAL && oprd2.GetType() != lexer.LITERAL {
			fmt.Printf("Erro: Operandos com tipos incompatíveis na linha %d, coluna %d. '%s' é do tipo '%s', enquanto que '%s' é do tipo '%s'\n", line, column, oprd1.GetLexem(), oprd1.GetType(), oprd2.GetLexem(), oprd2.GetType())
			semanticErrorFlag = true
			return
		}

		temporal := ""
		operationType := lexer.NULL

		switch oprd1.GetType() {
		case lexer.INTEGER:
			temporal = s.NewTemporal(TemporalInt)
			operationType = lexer.INTEGER
		case lexer.REAL:
			temporal = s.NewTemporal(TemporalFloat)
			operationType = lexer.REAL
		}

		s.AddToCodeBuffer(fmt.Sprintf("%s = %s %s %s;\n", temporal, oprd1.GetLexem(), opm.GetLexem(), oprd2.GetLexem()))
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), temporal, operationType)
		s.semanticStack.Push(newToken)
	},

	// LD -> OPRD
	20: func(s *Semantic, rule Rule, line int, column int) {
		oprdToken, _ := s.semanticStack.Pop()
		oprdTokenConverted := oprdToken.(lexer.Token)
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), oprdTokenConverted.GetLexem(), oprdTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// OPRD -> id
	21: func(s *Semantic, rule Rule, line int, column int) {
		idToken, _ := s.semanticStack.Pop()
		idTokenConverted := idToken.(lexer.Token)
		if idTokenConverted.GetType() == lexer.NULL {
			fmt.Printf("Erro: variável %s não declarada na linha %d, coluna %d\n", idTokenConverted.GetLexem(), line, column)
			semanticErrorFlag = true
			return
		}
		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), idTokenConverted.GetLexem(), idTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// OPRD -> num
	22: func(s *Semantic, rule Rule, line int, column int) {
		numToken, _ := s.semanticStack.Pop()
		numTokenConverted := numToken.(lexer.Token)

		newToken := lexer.NewToken(lexer.TokenClass(rule.Left), numTokenConverted.GetLexem(), numTokenConverted.GetType())
		s.semanticStack.Push(newToken)
	},

	// COND -> CAB CP
	24: func(s *Semantic, rule Rule, line int, column int) {
		s.AddToCodeBuffer("}\n")
	},

	// CAB -> se ab_p EXP_R fc_p entao
	25: func(s *Semantic, rule Rule, line int, column int) {
		s.semanticStack.Pop() // remove "entao" from stack
		s.semanticStack.Pop() // remove "fc_p" from stack
		rawExp_r, _ := s.semanticStack.Pop()
		exp_r := rawExp_r.(lexer.Token)
		s.AddToCodeBuffer(fmt.Sprintf("if (%s) {\n", exp_r.GetLexem()))
	},

	// EXP_R -> OPRD opr OPRD
	26: func(s *Semantic, rule Rule, line int, column int) {
		rawOprd2, _ := s.semanticStack.Pop()
		oprd2 := rawOprd2.(lexer.Token)

		rawOpr, _ := s.semanticStack.Pop()
		opr := rawOpr.(lexer.Token)

		rawOprd1, _ := s.semanticStack.Pop()
		oprd1 := rawOprd1.(lexer.Token)

		if oprd1.GetType() != oprd2.GetType() {
			fmt.Printf("Erro: Operandos com tipos incompatíveis na linha %d, coluna %d. '%s' é do tipo '%s', enquanto que '%s' é do tipo '%s'\n", line, column, oprd1.GetLexem(), oprd1.GetType(), oprd2.GetLexem(), oprd2.GetType())
			semanticErrorFlag = true
			return
		}

		temporalId := s.NewTemporal(TemporalBool)

		exp_rToken := lexer.NewToken(lexer.TokenClass(rule.Left), temporalId, lexer.NULL)
		s.semanticStack.Push(exp_rToken)
		s.AddToCodeBuffer(fmt.Sprintf("%s = %s %s %s;\n", temporalId, oprd1.GetLexem(), opr.GetLexem(), oprd2.GetLexem()))
	},

	// R -> CABR CPR
	32: func(s *Semantic, rule Rule, line int, column int) {
		s.AddToCodeBuffer("}\n")
	},

	// CABR -> repita ab_p EXP_R fc_p
	33: func(s *Semantic, rule Rule, line int, column int) {
		s.semanticStack.Pop() // remove "fc_p" from stack
		rawExp_r, _ := s.semanticStack.Pop()
		exp_r := rawExp_r.(lexer.Token)
		s.AddToCodeBuffer(fmt.Sprintf("while (%s) {\n", exp_r.GetLexem()))
	},
}

type Semantic struct {
	semanticStack *stack.Stack
	codeBuffer    *CodeBuffer
	ruleMap       map[int]func(s *Semantic, rule Rule, line int, column int)
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

func (s *Semantic) ExecuteRule(rule Rule, line int, column int) {
	_, found := s.ruleMap[rule.Number+1]
	if !found {
		return
	}
	s.ruleMap[rule.Number+1](s, rule, line, column)
}

func (s *Semantic) AddToCodeBuffer(code string) {
	s.codeBuffer.code += code
}

// NewTemporal adds a new temporal variable of TemporalType
func (s *Semantic) NewTemporal(temporalType TemporalType) string {
	temporalId := len(s.codeBuffer.temporals)
	s.codeBuffer.temporals = append(s.codeBuffer.temporals, temporalType)
	return fmt.Sprintf("T%d", temporalId)
}

func (s *Semantic) GenerateCode() {
	currentCode := `
#include<stdio.h>
#include<stdbool.h>
typedef char literal[256];
void main() {
`
	currentCode = fmt.Sprintf("%s%s", currentCode, s.codeBuffer.PrintTemporals())

	currentCode = fmt.Sprintf("%s%s", currentCode, s.codeBuffer.code)

	currentCode = fmt.Sprintf("%s%s", currentCode, "\n}")

	ioutil.WriteFile("output.c", []byte(currentCode), 0755)
}
