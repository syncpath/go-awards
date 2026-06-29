package parser_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/syncpath/go-awards/internal/parser"
)

func TestParseCSVFile(t *testing.T) {
	tests := []struct {
		in      string
		want    []map[string]string
		errWant bool
	}{
		{"Имя,Курс\nИванов,5\n", []map[string]string{{"Имя": "Иванов", "Курс": "5"}}, false},
		{"", nil, true},
		{"Имя,Курс\n", nil, false},
		{"Имя,Курс\nИванов\n", nil, true},
		{"Имя\n\"Иванов, И.И.\"\n", []map[string]string{{"Имя": "Иванов, И.И."}}, false},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case: %q, %v, %v", test.in, test.want, test.errWant)
		t.Run(name, func(t *testing.T) {
			got, errGot := parser.ParseCSVFile(strings.NewReader(test.in))
			if errGot == nil && test.errWant {
				t.Errorf("ожидал ошибку, получил nil")
			}
			if errGot != nil && !test.errWant {
				t.Errorf("ожидал nil, получил ошибку")
			}
			if !test.errWant && !reflect.DeepEqual(got, test.want) {
				t.Errorf("получил: %v, ожидал: %v", got, test.want)
			}
		})
	}
}
