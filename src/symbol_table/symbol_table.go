package symboltable

import (
	"fmt"
	"mgol-go/src/lexer"
)

// Posible errors
var (
	ErrorAlreadyOnTable = fmt.Errorf("the specified symbol is already on the symbol table")
	ErrorSymbolNotFound = fmt.Errorf("the specified symbol doesn't exists on the symbol table")
)

var table = map[string]lexer.Token{}

func Insert(id string, token lexer.Token) (lexer.Token, error) {
	_, found := table[id]
	if found {
		return lexer.Token{}, ErrorAlreadyOnTable
	}

	table[id] = token

	return table[id], nil
}

func GetToken(id string) (lexer.Token, error) {
	token, found := table[id]
	if !found {
		return lexer.Token{}, ErrorSymbolNotFound
	}
	return token, nil
}

func Update(id string, newToken lexer.Token) error {
	_, found := table[id]
	if !found {
		return ErrorSymbolNotFound
	}
	table[id] = newToken
	return nil
}

func CleanupTable() {
	for k := range table {
		delete(table, k)
	}
}

func GetTable() map[string]lexer.Token {
	return table
}
