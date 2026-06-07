// Package decliner склонение ФИО
package decliner

import (
	"errors"
	"fmt"
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

var petr, petrErr = Petrovich.LoadRules("Rules/rules.json")

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
		return "", errors.New("ФИО должно содержать минимум фамилию и имя")
	}

	val, ok := Cases[strings.ToLower(p)]
	if !ok {
		return "", errors.New("неверное название падежа")
	}

	// Petrovich принимает на вход в InfFio строку из 3 слов, если 2 слова, придется каждое по отдельности. Так как отчество неизвестно, то определить пол сложно. По умолчанию ставится мужской пол. Если человек оказался женского пола, то нужно вручную в итоговом .typ поменять склонение
	if len(fioWords) <= 2 {
		if !short {
			return strings.Join([]string{petr.InfLastname(fioWords[0], val, "male"), petr.InfFirstname(fioWords[1], val, "male")}, " "), nil
		} else {
			initial := string([]rune(fioWords[1])[0]) + "."
			return strings.Join([]string{petr.InfLastname(fioWords[0], val, "male"), initial}, " "), nil
		}
	}

	return petr.InfFio(fio, val, short), nil
}
