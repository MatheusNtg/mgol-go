package lexer

import (
	"fmt"
)

// Posible errors
var (
	ErrorAlreadyOnTable = fmt.Errorf("the specified symbol is already on the symbol table")
	ErrorSymbolNotFound = fmt.Errorf("the specified symbol doesn't exists on the symbol table")
)

var table = map[string]Token{}

func InsertSymbolTable(id string, token Token) (Token, error) {
	_, found := table[id]
	if found {
		return Token{}, ErrorAlreadyOnTable
	}

	table[id] = token

	return table[id], nil
}

func GetTokenFromSymbolTable(id string) (Token, error) {
	token, found := table[id]
	if !found {
		return Token{}, ErrorSymbolNotFound
	}
	return token, nil
}

func UpdateSymbolTable(id string, newToken Token) error {
	_, found := table[id]
	if !found {
		return ErrorSymbolNotFound
	}
	table[id] = newToken
	return nil
}

func CleanupSymbolTable() {
	for k := range table {
		delete(table, k)
	}
}

func GetSymbolTable() map[string]Token {
	return table
}

func PrintSymbolTable() {
	for k, v := range table {
		fmt.Printf("Chave: %v, Valor: %v\n", k, v)
	}
}
