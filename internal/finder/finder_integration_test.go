//go:build integration

package finder_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syncpath/go-awards/internal/finder"
	"github.com/syncpath/go-awards/internal/testutil"
)

func TestReadPdf(t *testing.T) {
	dir, err := testutil.FindPath("materials/examples")
	if err != nil {
		t.Fatalf("ошибка выполнения теста: %v", err)
	}

	tests := []struct {
		name     string
		dir      string
		file     string
		wantData finder.PageInfo
		errWant  bool
	}{
		{
			name: "корректный случай",
			dir:  dir,
			file: "hello.pdf",
			wantData: finder.PageInfo{
				Texts: []finder.TextField{
					{Text: "привет!!"},
				},
			},
			errWant: false,
		},
		{
			name:    "несуществующий файл",
			dir:     t.TempDir(),
			file:    "hello.pdf",
			errWant: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := finder.ReadPdf(test.dir, test.file)
			if gotErr == nil && test.errWant {
				t.Fatalf("ожидал ошибку, получил nil")
			}
			if gotErr != nil && !test.errWant {
				t.Fatalf("ожидал nil, получил ошибку: %v", gotErr)
			}

			var allText string
			for _, g := range got.Texts {
				allText += g.Text + " "
			}
			for _, text := range test.wantData.Texts {
				if !strings.Contains(allText, text.Text) {
					t.Fatalf("в полученном результате нет ожидаемого текста: %q", text.Text)
				}
			}
		})
	}

	// Тест с временным txt файлом
	t.Run("временный txt файл", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFileName := "broken.txt"
		file := filepath.Join(tmpDir, tmpFileName)
		err := os.WriteFile(file, []byte("это не pdf файл"), 0o644)
		if err != nil {
			t.Fatalf("не удалось создать временный файл: %v", err)
		}
		_, gotErr := finder.ReadPdf(tmpDir, tmpFileName)

		if gotErr == nil {
			t.Error("ожидал ошибку, получил nil")
		}
	})
}
