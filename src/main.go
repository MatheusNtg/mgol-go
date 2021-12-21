package main

import (
	"fmt"
	"io"
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

	buffer := make([]byte, 1)

	for {
		nBytes, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF {
			// TODO: Needs scanner to return tokens
			eofToken := lexer.NewToken(lexer.EOF, "", lexer.NULL)
			fmt.Print(eofToken)
			break
		}

		// TODO: We need to pass these results to the scanner
		fmt.Println("Reading: ", string(buffer[:nBytes]))
	}

}
