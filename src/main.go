package main

import (
	"log"
	"mgol-go/src/lexer"
	"mgol-go/src/parser"
	"mgol-go/src/stack"
	"os"
)

var separator = "=================="

const (
	stackCapacity   = 100000
	grammarPath     = "./src/parser/grammar.json"
	actionTablePath = "./src/parser/tables/action.tsv"
	gotoTablePath   = "./src/parser/tables/goto.tsv"
)

func main() {
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	symbolTable := lexer.GetSymbolTableInstance()

	lexer.FillSymbolTable(symbolTable)
	defer symbolTable.Cleanup()

	scanner := lexer.NewScanner(file, symbolTable)
	stack := stack.NewStack(stackCapacity)
	rules := parser.GetRulesMap(grammarPath)
	parser := parser.NewParser(scanner, stack, rules, actionTablePath, gotoTablePath)

	parser.Parse()
}
