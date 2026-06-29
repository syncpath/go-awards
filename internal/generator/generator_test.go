package generator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/syncpath/go-awards/internal/config"
)

func TestProcessTemplate(t *testing.T) {
	tests := []struct {
		value   string
		row     map[string]string
		want    string
		errWant bool
	}{
		{
			value:   "Привет, @val(\"name\")!",
			row:     map[string]string{"name": "Иван"},
			want:    "Привет, Иван!",
			errWant: false,
		},
		{
			value:   "Привет, @val(\"nonexistent\")!",
			row:     map[string]string{},
			want:    "Привет, !",
			errWant: false,
		},
		{
			value:   "@typst(\"text \\\"with\\\" quotes\")",
			row:     map[string]string{},
			want:    "text \"with\" quotes",
			errWant: false,
		},
		{
			value:   "Уважаемый(ая) @gender(\"name\", \"господин\", \"госпожа\")",
			row:     map[string]string{"name": "Иван"},
			want:    "Уважаемый(ая) господин",
			errWant: false,
		},
		{
			value:   "Уважаемый(ая) @gender(\"name\", \"господин\", \"госпожа\")",
			row:     map[string]string{"name": "Иванова Мария"},
			want:    "Уважаемый(ая) госпожа",
			errWant: false,
		},
		{
			value:   "Награждается @name(\"name\", дательный, false)",
			row:     map[string]string{"name": "Иванов Иван"},
			want:    "Награждается Иванову Ивану",
			errWant: false,
		},
		{
			value:   "Награждается @name(\"name\", дательный, false, 1)",
			row:     map[string]string{"name": "Иванов Иван"},
			want:    "Награждается Иванову \\ Ивану",
			errWant: false,
		},
		{
			value:   "Награждается @name(\"name\", дательный, false, 0)",
			row:     map[string]string{"name": "Иванов Иван"},
			want:    "",
			errWant: true,
		},
		{
			value:   "Награждается @name(\"name\", дательный, false, 2)",
			row:     map[string]string{"name": "Иванов Иван"},
			want:    "",
			errWant: true,
		},
		{
			value:   "Награждается @name(\"name\", дательный, false, 3)",
			row:     map[string]string{"name": "Иванов Иван"},
			want:    "",
			errWant: true,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case: %s", test.value)
		t.Run(name, func(t *testing.T) {
			got, errGot := processTemplate(test.value, test.row)
			if errGot == nil && test.errWant {
				t.Errorf("ожидал ошибку, получил nil")
			}
			if errGot != nil && !test.errWant {
				t.Errorf("ожидал nil, получил ошибку: %v", errGot)
			}
			if got != test.want {
				t.Errorf("получил: %s, ожидал: %s", got, test.want)
			}
		})
	}
}

func TestCycleCheck(t *testing.T) {
	tests := []struct {
		name             string
		field            config.Field
		parentToChildren map[string][]config.Field
		errWant          bool
	}{
		{
			name:  "Линейное дерево без циклов",
			field: config.Field{ID: "A"},
			parentToChildren: map[string][]config.Field{
				"A": {{ID: "B"}},
				"B": {{ID: "C"}},
			},
			errWant: false,
		},
		{
			name:  "Явный простой цикл",
			field: config.Field{ID: "A"},
			parentToChildren: map[string][]config.Field{
				"A": {{ID: "B"}},
				"B": {{ID: "A"}},
			},
			errWant: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			visited := make(map[string]bool)
			visiting := make(map[string]bool)
			errGot := CycleCheck(test.field, test.parentToChildren, visited, visiting)
			if errGot == nil && test.errWant {
				t.Errorf("ожидал ошибку, получил nil")
			}
			if errGot != nil && !test.errWant {
				t.Errorf("ожидал nil, получил ошибку: %v", errGot)
			}
		})
	}
}

func TestGenerateTypst(t *testing.T) {
	cfg := &config.TemplateConfig{
		Width:  210.0,
		Height: 297.0,
		Bg: config.Background{
			BackgroundImage: "bg.png",
			BgScale:         100.0,
		},
		Fields: []config.Field{
			{
				ID:    "title",
				Value: "Диплом для @val(\"name\")",
				X:     10.0,
				Y:     20.0,
				Width: 100.0,
			},
		},
	}

	table := []map[string]string{
		{"name": "Ивана"},
	}

	res, errGot := GenerateTypst(cfg, table)
	if errGot != nil {
		t.Fatalf("ожидал nil, получил ошибку: %v", errGot)
	}

	if !strings.Contains(res, "Диплом для Ивана") {
		t.Errorf("получил: %s, ожидал вхождение: 'Диплом для Ивана'", res)
	}
}
