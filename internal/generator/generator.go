// Package generator Генерация .typ файла
package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/decliner"
)

var (
	reName = regexp.MustCompile(`@name\("([^"]+)",\s*([^,]+),\s*(true|false)\)`)
	reVal  = regexp.MustCompile(`@val\("([^"]+)"\)`)
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
)
`,
		cfg.Width,
		cfg.Height,
		cfg.Bg.BgDx,
		cfg.Bg.BgDy,
		cfg.Bg.BackgroundImage,
		cfg.Bg.BgScale,
		cfg.Bg.BgScale)

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

	// А теперь для каждого участника
	for i, row := range table {
		// Разрыв страницы для каждого участника кроме первого
		if i > 0 {
			sb.WriteString("\n#pagebreak()\n\n")
		}

		// Расставляем поля для конкретного участника
		for _, rootField := range rootFields {
			rootValue, err := processTemplate(rootField.Value, row)
			if err != nil {
				return "", fmt.Errorf("ошибка замены плейсхолдера: %w", err)
			}

			fmt.Fprintf(&sb, "#place(top + left, dx: %.3fmm, dy: %.3fmm)[\n", rootField.X, rootField.Y)

			fmt.Fprintf(&sb,
				`  #block(width: %.3fmm)[
		#set text(font: "%s", size: %.3fpt, weight: "%s")
		#align(%s)[%s]
	]
`, rootField.Width, rootField.Font, rootField.FontSize, rootField.FontType, rootField.Align, rootValue)

			for _, child := range parentToChildren[rootField.ID] {
				childValue, err := processTemplate(child.Value, row)
				if err != nil {
					return "", fmt.Errorf("ошибка замены плейсхолдера: %w", err)
				}
				fmt.Fprintf(&sb,
					`	 #v(%.3fmm)  
						 #block(width: %.3fmm)[
						   #set text(font: "%s", size: %.3fpt, weight: "%s")
							 #align(%s)[%s]
						 ]
					`, child.IndentValue, child.Width, child.Font, child.FontSize, child.FontType, child.Align, childValue)
			}
			sb.WriteString("]\n")
		}
	}
	return sb.String(), nil
}

// processTemplate поиск шаблонов @name(...) и @val(...) и замена их на значения из csv. У name происходит дополнительно склонение
func processTemplate(value string, row map[string]string) (string, error) {
	result := value

	// Для @name(...)
	matches := reName.FindAllStringSubmatch(result, -1)
	for _, match := range matches {
		fullMacro := match[0]
		fio := row[match[1]]
		p := match[2]
		isShort := match[3] == "true"

		r, err := decliner.Decline(fio, p, isShort)
		if err != nil {
			return "", fmt.Errorf("не удалось склонить ФИО: %w", err)
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

	return result, nil
}
