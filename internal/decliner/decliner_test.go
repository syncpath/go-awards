package decliner_test

import (
	"fmt"
	"testing"

	"github.com/syncpath/go-awards/internal/decliner"
)

func TestDetectGender(t *testing.T) {
	tests := []struct {
		fio  string
		want string
	}{
		{"Иванов Пётр Васильевич", "male"},
		{"Иванова Екатерина Дмитриевна", "female"},
		{"Дрозодов Кирилл", "male"},
		{"Петрова Дарья", "female"},
		{"Абдуллаев Руслан Ильгар оглы", "male"},
		{"Кирилл", "male"},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case: %s", test.fio)
		t.Run(name, func(t *testing.T) {
			got := decliner.DetectGender(test.fio)
			if got != test.want {
				t.Errorf("получил: %s, ожидал: %s", got, test.want)
			}
		})
	}
}

func TestDecline(t *testing.T) {
	tests := []struct {
		fio     string
		p       string
		short   bool
		want    string
		errWant bool
	}{
		{"Иванов Пётр Васильевич", "ИмеНиТельНый", false, "Иванов Пётр Васильевич", false},
		{"Петров Иван Кириллович", "именительный", true, "Петров И.К.", false},
		{"Дрозденко Евгений Иванович", "Родительный", false, "Дрозденко Евгения Ивановича", false},
		{"Столл Елена Анатольевна", "дАтельный", false, "Столл Елене Анатольевне", false},
		{"Белькович Егор Алексеевич", "Творительный", false, "Бельковичем Егором Алексеевичем", false},
		{"Белькович Егор Алексеевич", "Творител", false, "", true},
		{"Зайдуллин Иван", "Винительный", true, "Зайдуллина И.", false},
		{"Петров Иван", "Родительный", false, "Петрова Ивана", false},
		{"Белькович", "Творительный", true, "", true},
		{"Абдуллаев Руслан Ильгар оглы", "Творительный", true, "", true},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case: %s, %v, %v", test.fio, test.p, test.short)
		t.Run(name, func(t *testing.T) {
			got, errGot := decliner.Decline(test.fio, test.p, test.short)
			if errGot == nil && test.errWant {
				t.Errorf("ожидал ошибку, получил nil")
			}
			if errGot != nil && !test.errWant {
				t.Errorf("ожидал nil, получил ошибку")
			}
			if got != test.want {
				t.Errorf("получил: %s, ожидал: %s", got, test.want)
			}
		})
	}
}
