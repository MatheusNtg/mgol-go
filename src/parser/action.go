package parser

import (
	"encoding/csv"
	"mgol-go/src/lexer"
	"os"
)

type Action int

const (
	SHIFT Action = iota
	REDUCE
)

type ActionReader struct {
	records [][]string
	indexes map[string]int
}

func NewActionReader(path string) *ActionReader {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	ac := &ActionReader{}
	ac.indexes = make(map[string]int)
	for idx, record := range records[0] {
		ac.indexes[record] = idx
	}
	ac.records = records
	return ac
}

func (a *ActionReader) GetAction(state lexer.State, token string) string {
	//We need to sum to sum one to access line n because we want to
	//eliminate the header itself
	return a.records[state+1][a.indexes[token]]
}
