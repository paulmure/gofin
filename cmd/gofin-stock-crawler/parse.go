package main

import (
	"encoding/csv"
	"io"
	"os"
)

func parseFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	res := make([]string, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, record[0])
	}

	return res, nil
}
