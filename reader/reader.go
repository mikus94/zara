package reader

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCSVFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file, due: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read csv lines, due: %w", err)
	}
	return lines, nil
}
