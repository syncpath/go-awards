// Package parser работа с CSV таблицами
package parser

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// ParseCSV загрузка и обработка CSV таблицы
func ParseCSV(dir, name string) ([]map[string]string, error) {
	fullpath := filepath.Join(dir, name)
	file, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseCSVFile(file)
}

// ParseCSVFile обработка CSV таблицы
func ParseCSVFile(r io.Reader) ([]map[string]string, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 || len(records[0]) == 0 {
		return nil, errors.New("пустой csv файл")
	}

	var data []map[string]string

	for i := 1; i < len(records); i++ {
		m := make(map[string]string)

		for j := 0; j < len(records[0]); j++ {
			m[records[0][j]] = records[i][j]
		}
		data = append(data, m)
	}

	return data, nil
}
