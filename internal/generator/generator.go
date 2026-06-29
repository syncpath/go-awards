// Package generator Генерация .typ файла
package generator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/decliner"
)

var (
	reName   = regexp.MustCompile(`@name\("([^"]+)",\s*([^,]+),\s*(true|false)(?:,\s*(\d+))?\)`)
	reVal    = regexp.MustCompile(`@val\("([^"]+)"\)`)
	reTypst  = regexp.MustCompile(`@typst\("((?:[^"\\]|\\.)*)"\)`)
	reGender = regexp.MustCompile(`@gender\("([^"]+)",\s*"([^"]+)",\s*"([^"]+)"\)`)
)

func GenerateTypst(cfg *config.TemplateConfig, table []map[string]string) (string, error) {
	var sb strings.Builder

	// Настройка страницы
	fmt.Fprintf(&sb,
		`#set page(
	width: %.3fmm,
	height: %.3fmm,
	margin: 0pt,
	background: place(
		center + horizon,
		dx: %.3fmm,
		dy: %.3fmm,
		image("%s", width: %.3f%%, height: %.3f%%),
	),
)`,
		cfg.Width,
		cfg.Height,
		cfg.Bg.BgDx,
		cfg.Bg.BgDy,
		cfg.Bg.BackgroundImage,
		cfg.Bg.BgScale,
		cfg.Bg.BgScale)

	sb.WriteString("\n#set block(spacing: 0pt)\n")
	sb.WriteString("\n")

	// Карты для быстрого поиска связей. Создаем мапу со связями. Если Indent не пустой, то создаем сына в родительском #place, если пустой, то в отдельном)
	parentToChildren := make(map[string][]config.Field)
	var rootFields []config.Field

	for _, field := range cfg.Fields {
		if field.Indent == "" {
			rootFields = append(rootFields, field)
		} else {
			parentToChildren[field.Indent] = append(parentToChildren[field.Indent], field)
		}
	}

	// Проверяем всю конфигурацию на циклы
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	for _, field := range cfg.Fields {
		if err := CycleCheck(field, parentToChildren, visited, visiting); err != nil {
			return "", fmt.Errorf("ошибка валидации шаблона: %w", err)
		}
	}

	// Проверяем поля ID json на повторы
	check := make(map[string]bool)
	for _, field := range cfg.Fields {
		if check[field.ID] {
			return "", errors.New("ошибка валидации шаблона: найдены повторы ID полей")
		}
		check[field.ID] = true
	}

	// А теперь для каждого участника
	for i, row := range table {
		// Разрыв страницы для каждого участника кроме первого
		if i > 0 {
			sb.WriteString("\n#pagebreak()\n\n")
		}

		fmt.Fprintf(&sb, "// Страница №%v \n", i+1)

		// Расставляем поля для конкретного участника
		for _, rootField := range rootFields {
			if rootField.PlaceAlign == "" {
				rootField.PlaceAlign = "top + left"
			}

			fmt.Fprintf(&sb, "#place(%s, dx: %.3fmm, dy: %.3fmm)[\n", rootField.PlaceAlign, rootField.X, rootField.Y)

			err := RenderField(rootField, parentToChildren, row, &sb, 1)
			if err != nil {
				return "", fmt.Errorf("ошибка генерации typst: %w", err)
			}

			sb.WriteString("]\n\n")
		}
	}
	return sb.String(), nil
}

func CycleCheck(field config.Field, parentToChild map[string][]config.Field, visited map[string]bool, visiting map[string]bool) error {
	if visiting[field.ID] {
		return errors.New("обнаружена циклическая зависимость или два поля указывают на одного ребенка")
	}
	if visited[field.ID] {
		return nil
	}

	visiting[field.ID] = true

	for _, child := range parentToChild[field.ID] {
		if err := CycleCheck(child, parentToChild, visited, visiting); err != nil {
			return err
		}
	}

	delete(visiting, field.ID)
	visited[field.ID] = true
	return nil
}

// RenderField вспомогательная функция для рекурсивного обхода
func RenderField(field config.Field, parentToChild map[string][]config.Field, row map[string]string, sb *strings.Builder, depth int) error {
	indent := strings.Repeat("  ", depth)
	value, err := processTemplate(field.Value, row)
	if err != nil {
		return fmt.Errorf("ошибка замены плейсхолдера: %w", err)
	}

	if field.Color == "" {
		field.Color = "#000000"
	}
	if field.Spacing == 0.0 {
		field.Spacing = 100.0
	}
	fmt.Fprintf(sb,
		"%[1]s#block(width: %.3[2]fmm)[\n"+
			"%[1]s  #set text(font: \"%[3]s\", size: %.3[4]fpt, tracking: %.3[10]fpt, spacing: %.3[11]f%%, weight: \"%[5]s\", fill: rgb(\"%[6]s\"))\n"+
			"%[1]s  #set par(leading: %.3[7]fmm)\n"+
			"%[1]s  #align(%[8]s)[%[9]s]\n"+
			"%[1]s]\n",
		indent, field.Width, field.Font, field.FontSize, field.FontType, field.Color, field.Leading, field.Align, value, field.Tracking, field.Spacing)

	for _, child := range parentToChild[field.ID] {
		fmt.Fprintf(sb, "%s#v(%.3fmm)\n", indent+"  ", child.IndentValue)

		if err := RenderField(child, parentToChild, row, sb, depth+1); err != nil {
			return err
		}
	}
	return nil
}

// processTemplate поиск шаблонов @name(...) и @val(...) и замена их на значения из csv. У name происходит дополнительно склонение
func processTemplate(value string, row map[string]string) (string, error) {
	result := value

	// Для @typst(...)
	matches := reTypst.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		fullMacro := match[0]
		cleaned := strings.ReplaceAll(match[1], "\\\"", "\"")
		result = strings.ReplaceAll(result, fullMacro, cleaned)
	}

	// Для @name(...)
	matches = reName.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		fullMacro := match[0]
		fio := row[match[1]]
		p := match[2]
		isShort := match[3] == "true"
		var lineBr int
		var err error

		r, err := decliner.Decline(fio, p, isShort)
		if err != nil {
			return "", fmt.Errorf("не удалось склонить ФИО: %w", err)
		}

		if match[4] != "" {
			lineBr, err = strconv.Atoi(match[4])
			if err != nil {
				return "", fmt.Errorf("ошибка для идентификатора переноса: %w", err)
			}

			sl := strings.Fields(r)
			if lineBr <= 0 || lineBr >= len(sl) {
				return "", errors.New("идентификатор переноса должен быть больше нуля, но меньше двух")
			}

			part1 := strings.Join(sl[:lineBr], " ")
			part2 := strings.Join(sl[lineBr:], " ")

			r = part1 + " \\ " + part2
		}

		result = strings.ReplaceAll(result, fullMacro, r)
	}

	// Для @val(...)
	matches = reVal.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		fullMacro := match[0]
		str := row[match[1]]

		result = strings.ReplaceAll(result, fullMacro, str)
	}

	// Для @gender(...)
	matches = reGender.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		fullMacro := match[0]
		if decliner.DetectGender(row[match[1]]) == "male" {
			result = strings.ReplaceAll(result, fullMacro, match[2])
		} else {
			result = strings.ReplaceAll(result, fullMacro, match[3])
		}
	}

	return result, nil
}
