//go:build integration

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/testutil"
)

func TestLoadConfigRealFS(t *testing.T) {
	dir, err := testutil.FindPath("templates")
	if err != nil {
		t.Fatalf("папка templates не найдена")
	}

	tests := []struct {
		name    string
		dir     string
		file    string
		wantErr bool
	}{
		{
			name: "корректный случай",
			dir:  dir,
			file: "hello.json",
		},
		{
			name:    "несуществующий файл",
			dir:     t.TempDir(),
			file:    "hello.json",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := config.LoadConfig(config.RealFS{}, test.dir, test.file)
			if gotErr == nil && test.wantErr {
				t.Fatalf("ожидал ошибку, получил nil")
			}
			if gotErr != nil && !test.wantErr {
				t.Fatalf("ожидал nil, получил ошибку: %v", gotErr)
			}
			if test.wantErr {
				return
			}

			if got.Title != "hello" {
				t.Errorf("получил %q, ожидал %q", got.Title, "hello")
			}
			if got.Width <= 0 || got.Height <= 0 {
				t.Errorf("Height: получил Width=%v Height=%v, ожидал > 0", got.Width, got.Height)
			}
			if len(got.Fields) != 1 {
				t.Fatalf("Fields: получил %d, ожидал 1", len(got.Fields))
			}
			if got.Fields[0].ID != "name" {
				t.Errorf("Fields[0].ID: получил %q, ожидал %q", got.Fields[0].ID, "name")
			}
		})
	}

	t.Run("битый json", func(t *testing.T) {
		tmpDir := t.TempDir()
		name := "broken.json"
		fullpath := filepath.Join(tmpDir, name)

		err := os.WriteFile(fullpath, []byte(`{"title": "битый`), 0o644)
		if err != nil {
			t.Fatalf("не удалось создать временный файл: %v", err)
		}
		_, gotErr := config.LoadConfig(config.RealFS{}, tmpDir, name)

		if gotErr == nil {
			t.Error("ожидал ошибку, получил nil")
		}
	})
}
