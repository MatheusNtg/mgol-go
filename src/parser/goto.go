package parser

import (
	"encoding/csv"
	"os"
)

type GotoReader struct {
	records [][]string
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

	return &GotoReader{
		records: records,
	}
}
