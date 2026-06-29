package config_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/syncpath/go-awards/internal/config"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		in      string
		want    *config.TemplateConfig
		errWant bool
	}{
		{`{
	"title": "Full Config",
	"width": 200,
	"background": {
		"background_image": "bg.png",
		"bg_scale": 105.5
	},
	"fields": [
		{"id": "name", "x": 10, "y": 20, "font_size": 12}
	]
}`, &config.TemplateConfig{"Full Config", 200.0, 0.0, []config.Field{{ID: "name", X: 10.0, Y: 20.0, FontSize: 12}}, config.Background{BackgroundImage: "bg.png", BgScale: 105.5}}, false},
		{`{
	"title": "Test Certificate",
	"width": 297.0
}`, &config.TemplateConfig{"Test Certificate", 297.0, 0.0, nil, config.Background{}}, false},
		{`{"title": "Bad JSON",`, nil, true},
		{`{"title": "Bad Type", "width": "200"}`, nil, true},
		{`{}`, &config.TemplateConfig{}, false},
	}
	for _, test := range tests {
		name := fmt.Sprintf("case: %q, %v, %v", test.in, test.want, test.errWant)
		t.Run(name, func(t *testing.T) {
			got, errGot := config.ParseConfig([]byte(test.in))
			if errGot == nil && test.errWant {
				t.Errorf("получил nil, ожидал err")
			}
			if errGot != nil && !test.errWant {
				t.Errorf("получил err, ожидал nil")
			}
			if !test.errWant && !reflect.DeepEqual(got, test.want) {
				t.Errorf("получил: %v, ожидал: %v", got, test.want)
			}
		})
	}
}

type MockFS struct {
	bytes []byte
	err   error
}

func (m MockFS) ReadFile(name string) ([]byte, error) {
	return m.bytes, m.err
}

func TestLoadConfig(t *testing.T) {
	t.Run("успешное чтение файла", func(t *testing.T) {
		mock := MockFS{
			bytes: []byte(`{"title": "hello from the mock"}`),
			err:   nil,
		}

		got, err := config.LoadConfig(mock, "dummy.json")
		if err != nil {
			t.Fatalf("получил err, ожидал nil")
		}
		if got.Title != "hello from the mock" {
			t.Errorf("получил %q, ожидал: hello from the mock", got.Title)
		}
	})

	t.Run("чтение несуществующего файла", func(t *testing.T) {
		mock := MockFS{
			bytes: nil,
			err:   os.ErrNotExist,
		}
		_, err := config.LoadConfig(mock, "missing.json")
		if err == nil {
			t.Fatalf("получил nil, ожидал err")
		}
	})
}
