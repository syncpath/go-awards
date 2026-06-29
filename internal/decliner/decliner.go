// Package decliner склонение ФИО
package decliner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fyrolab/Petrovich-Go"
)

// Cases Мапа со всеми падежами Petrovich (в библиотеки они оформлены с помощью iota)
var Cases = map[string]int{
	"родительный":  Petrovich.Genitive,      // 0
	"дательный":    Petrovich.Dative,        // 1
	"винительный":  Petrovich.Accusative,    // 2
	"творительный": Petrovich.Instrumental,  // 3
	"предложный":   Petrovich.Prepositional, // 4
}

func findRulesPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		path := filepath.Join(dir, "Rules", "rules.json")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("Rules/rules.json не найдена")
}

var (
	petrPath, _   = findRulesPath()
	petr, petrErr = Petrovich.LoadRules(petrPath)
)

func InitDecliner(rulesPath string) error {
	if _, err := os.Stat(rulesPath); err != nil && !filepath.IsAbs(rulesPath) {
		if path, walkErr := findRulesPath(); walkErr == nil {
			rulesPath = path
		}
	}
	petr, petrErr = Petrovich.LoadRules(rulesPath)
	return petrErr
}

/*
Decline Перевод в нужный падеж + сокращение до инициалов при необходимости.
------------------------------------------------------------------------------------
fio - ФИО, p - нужный падеж на русском языке, short - сокращать до инициалов или нет
*/
func Decline(fio string, p string, short bool) (string, error) {
	if petrErr != nil {
		return "", fmt.Errorf("ошибка инициализации правил Petrovich: %w", petrErr)
	}

	fioWords := strings.Fields(fio)
	l := len(fioWords)

	if l < 2 {
		return "", errors.New("ФИО должно содержать хотя бы фамилию и имя")
	}

	if l > 3 {
		return "", errors.New("ФИО должно быть не более 3 слов. Eсли есть ФИО с 4 и более словами, в таблице оставьте только 3, в итоговом typst исправьте ФИО для данного человека")
	}

	if strings.ToLower(p) == "именительный" {
		if short {
			initial := string([]rune(fioWords[1])[0]) + "."
			if len(fioWords) >= 3 {
				initial += string([]rune(fioWords[2])[0]) + "."
			}
			return strings.Join([]string{fioWords[0], initial}, " "), nil
		}
		return strings.Join(fioWords, " "), nil
	}

	val, ok := Cases[strings.ToLower(p)]
	if !ok {
		return "", errors.New("неверное название падежа")
	}

	// Petrovich принимает на вход в InfFio строку из 3 слов, если 2 слова, придется каждое по отдельности. Так как отчество неизвестно, то определить пол сложно. По умолчанию ставится мужской пол. Если человек оказался женского пола, то нужно вручную в итоговом .typ поменять склонение
	if len(fioWords) <= 2 {

		gender := DetectGender(fio)
		if !short {
			fioOut := strings.Join([]string{petr.InfLastname(fioWords[0], val, gender), petr.InfFirstname(fioWords[1], val, gender)}, " ")
			return fioOut, nil
		} else {
			initial := string([]rune(fioWords[1])[0]) + "."
			fioOut := strings.Join([]string{petr.InfLastname(fioWords[0], val, gender), initial}, " ")
			return fioOut, nil
		}
	}

	return petr.InfFio(fio, val, short), nil
}

func DetectGender(fio string) string {
	words := strings.Fields(fio)
	if len(words) < 3 {
		for _, word := range words {
			w := strings.ToLower(word)

			if strings.HasSuffix(w, "ова") || strings.HasSuffix(w, "ева") || strings.HasSuffix(w, "ина") || strings.HasSuffix(w, "ая") {
				return "female"
			}
		}
		return "male"
	}

	patronymic := strings.ToLower(words[2])

	if strings.HasSuffix(patronymic, "ич") {
		return "male"
	} else if strings.HasSuffix(patronymic, "на") {
		return "female"
	}
	return "male"
}
