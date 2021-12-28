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

	// Temporary implementation, before the pattern recognizer
	for {
		currToken := lexer.Scanner(file)
		fmt.Println(currToken)

		if currToken == lexer.EOF_TOKEN {
			break
		}
	}

}
