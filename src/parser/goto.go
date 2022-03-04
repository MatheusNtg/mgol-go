package parser

import (
	"encoding/csv"
	"mgol-go/src/lexer"
	"os"
	"strconv"
)

type GotoReader struct {
	records [][]string
	indexes map[string]int
}

func NewGotoReader(path string) *GotoReader {
	gotoCsvFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer gotoCsvFile.Close()

	gotoCsv := csv.NewReader(gotoCsvFile)
	gotoCsv.Comma = '\t'

	records, err := gotoCsv.ReadAll()
	if err != nil {
		panic(err)
	}

	got := &GotoReader{}
	got.indexes = make(map[string]int)
	for idx, record := range records[0] {
		got.indexes[record] = idx
	}
	got.records = records
	return got
}

func (g *GotoReader) GetGoto(state lexer.State, nonTerminal string) int {
	//We need to sum to sum one to access line n because we want to
	//eliminate the header itself
	value := g.records[state+1][g.indexes[nonTerminal]]

	if len(value) == 0 {
		return -1
	}

	ret_state, err := strconv.Atoi(string(value))
	if err != nil {
		panic(err)
	}

	return ret_state
}
