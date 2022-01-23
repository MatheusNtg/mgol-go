package main

import (
	"fmt"
	"log"
	"mgol-go/src/lexer"
	"os"
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
	for {
		token := scanner.Scan()
		fmt.Println(token)
		if token == lexer.EOF_TOKEN {
			break
		}
	}
}
