package lexer

import (
	"io"
	"log"
	"os"
)

// Scans a code file
func Scanner(file *os.File) Token {
	buffer := make([]byte, 1)
	_, err := file.Read(buffer)

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	if err == io.EOF {
		return EOF_TOKEN
	}

	return Token{}
}
