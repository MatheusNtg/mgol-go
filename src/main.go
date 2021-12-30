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

	scanner := lexer.NewScanner(file)

	for {
		token := scanner.Scan()
		fmt.Println(token)
		if token == lexer.EOF_TOKEN {
			return
		}
	}

}
