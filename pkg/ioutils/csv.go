package ioutils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCsvFile tries to open and read a CSV file on the given path as a slice of string slices.
func ReadCsvFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file \"%s\"", path)
	}

	defer func(f *os.File) {
		cerr := f.Close()
		if cerr != nil {
			err = cerr
		}
	}(f)

	csvReader := csv.NewReader(f)
	return csvReader.ReadAll()
}
