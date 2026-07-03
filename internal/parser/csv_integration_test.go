//go:build integration

package parser_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/syncpath/go-awards/internal/parser"
	"github.com/syncpath/go-awards/internal/testutil"
)

func TestParseCSV(t *testing.T) {
	dir, err := testutil.FindPath("materials/tables")
	if err != nil {
		t.Fatalf("папка tables не найдена")
	}

	tests := []struct {
		name    string
		dir     string
		file    string
		want    []map[string]string
		errWant bool
	}{
		{
			name:    "корректный случай",
			dir:     dir,
			file:    "hello.csv",
			want:    []map[string]string{{"привет": "тест!"}},
			errWant: false,
		},
		{
			name:    "несуществующий_файл",
			dir:     t.TempDir(),
			file:    "lol.csv",
			errWant: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := parser.ParseCSV(test.dir, test.file)
			if gotErr == nil && test.errWant {
				t.Fatalf("ожидал ошибку, получил nil")
			}
			if gotErr != nil && !test.errWant {
				t.Fatalf("ожидал nil, получил ошибку: %v", gotErr)
			}
			if !reflect.DeepEqual(got, test.want) && !test.errWant {
				t.Errorf("ожидал: %q, получил: %q", test.want, got)
			}
		})
	}

	// Тест файла текстового временного
	t.Run("кривой CSV", func(t *testing.T) {
		dir := t.TempDir()
		name := "hello.txt"
		fullpath := filepath.Join(dir, name)

		err := os.WriteFile(fullpath, []byte("Имя,Курс\nИванов\n"), 0o644)
		if err != nil {
			t.Fatalf("не удалось создать временный файл: %v", err)
		}
		_, gotErr := parser.ParseCSV(dir, name)

		if gotErr == nil {
			t.Error("ожидал ошибку, получил nil")
		}
	})
}
