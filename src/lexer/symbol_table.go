package lexer

import (
	"fmt"
)

// Posible errors
var (
	ErrorAlreadyOnTable = fmt.Errorf("the specified symbol is already on the symbol table")
	ErrorSymbolNotFound = fmt.Errorf("the specified symbol doesn't exists on the symbol table")
)

type SymbolTable struct {
	table map[string]Token
}

var symbolTableInstance *SymbolTable

func GetSymbolTableInstance() *SymbolTable {
	if symbolTableInstance == nil {
		symbolTableInstance = &SymbolTable{
			table: make(map[string]Token),
		}
		return symbolTableInstance
	}
	return symbolTableInstance
}

func (s *SymbolTable) Insert(id string, token Token) Token {
	tok, found := s.table[id]
	if found {
		return tok
	}

	s.table[id] = token

	return s.table[id]
}

func (s *SymbolTable) GetToken(lexem string) (Token, error) {
	token, found := s.table[lexem]
	if !found {
		return Token{}, ErrorSymbolNotFound
	}
	return token, nil
}

func (s *SymbolTable) Update(id string, newToken Token) error {
	_, found := s.table[id]
	if !found {
		return ErrorSymbolNotFound
	}
	s.table[id] = newToken
	return nil
}

func (s *SymbolTable) Cleanup() {
	for k := range s.table {
		delete(s.table, k)
	}
}

func (s *SymbolTable) Print() {
	for k, v := range s.table {
		fmt.Printf("Chave: %v, Valor: %v\n", k, v)
	}
}
