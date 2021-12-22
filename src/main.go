package main

import (
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

	// Temporary implementation, before the pattern recognizer
	for {
		currToken := lexer.Scanner(file)
		if currToken == *lexer.NewToken(lexer.EOF, "", lexer.NULL) {
			break
		}
	}
}
