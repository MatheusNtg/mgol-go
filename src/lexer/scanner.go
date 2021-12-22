package lexer

import (
	"io"
	"log"
	"os"
)

func Scanner(file *os.File) *Token {
	buffer := make([]byte, 1)
	nBytes, err := file.Read(buffer)

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	if err == io.EOF {
		EofToken := NewToken(EOF, "", NULL)
		return EofToken
	}
}
