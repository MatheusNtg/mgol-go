package parser

import (
	"encoding/csv"
	"mgol-go/src/lexer"
	"os"
	"strconv"
	"strings"
)

type Action int

const (
	SHIFT Action = iota
	REDUCE
	ACCEPT
	ERROR
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

func (a *ActionReader) GetAction(state lexer.State, token lexer.Token) (Action, int) {
	var class string
	if token == lexer.EOF_TOKEN {
		class = "$"
	} else {
		class = token.GetClass()
	}

	//We need to sum to sum one to access line n because we want to
	//eliminate the header itself
	value := []byte(a.records[state+1][a.indexes[class]])

	if len(value) == 0 {
		return ERROR, 0
	}

	if strings.Compare(string(value), "acc") == 0 {
		return ACCEPT, 0
	}

	switch value[0] {
	case 's':
		opr, err := strconv.Atoi(string(value[1:]))
		if err != nil {
			panic(err)
		}
		return SHIFT, opr
	case 'r':
		opr, err := strconv.Atoi(string(value[1:]))
		if err != nil {
			panic(err)
		}
		return REDUCE, opr
	case 'e':
		errorType, err := strconv.Atoi(string(value[1:]))
		if err != nil {
			panic(err)
		}
		return ERROR, errorType
	}

	return -1, -1
}
