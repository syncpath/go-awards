// Package finder Изучает наш макет (например, из брендбука) и выводит всю информацию о нем
package finder

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

var exFolder string = "examples/"

const coefmm float64 = 25.4 / 72.0

// MergedText вспомогательная структура, нужна, так как ledongthuc/pdf хранит каждый символ как отдельный текст
type MergedText struct {
	Font     string
	FontSize float64
	X        float64
	Y        float64
	S        string
	W        float64
}

// FindFields показывает координаты текстовых полей, которые она смогла найти
func FindFields(path string) ([]string, error) {
	fullpath := filepath.Join(exFolder, path)
	file, reader, err := pdf.Open(fullpath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения макета: %w", err)
	}

	defer file.Close()

	var res []string

	// Подразумевается, что макет имеет только 1 страницу
	p := reader.Page(1)

	if p.V.IsNull() {
		return nil, errors.New("страница пустая")
	}
	coefsMap := FindFontsCoefs(p)
	texts := p.Content().Text

	var merged []MergedText

	for _, t := range texts {
		if t.S == "" {
			continue
		}

		if len(merged) > 0 {
			last := &merged[len(merged)-1]

			if last.Font == t.Font && last.FontSize == t.FontSize && last.Y == t.Y {
				last.S += t.S
				continue
			}
		}

		merged = append(merged, MergedText{
			Font:     t.Font,
			FontSize: t.FontSize,
			X:        t.X,
			Y:        t.Y,
			S:        t.S,
		})
	}

	// Параметры страницы
	mediaBox := p.V.Key("MediaBox")
	xMin, xMax := mediaBox.Index(0).Float64(), mediaBox.Index(2).Float64()
	yMin, yMax := mediaBox.Index(1).Float64(), mediaBox.Index(3).Float64()

	height := (yMax - yMin) * coefmm
	width := (xMax - xMin) * coefmm
	pageStr := fmt.Sprintf("Page: height: %.1fmm, width: %1.fmm", height, width)

	res = append(res, pageStr)
	for _, text := range merged {
		cleanName := cleanFontName(text.Font)
		coef, ok := coefsMap[cleanName]
		if !ok {
			coef = 0.7
		}

		x := (text.X - xMin) * coefmm
		y := (yMax - (text.Y + text.FontSize*coef)) * coefmm

		runes := []rune(text.S)
		var displayText string
		if len(runes) > 15 {
			displayText = string(runes[:15]) + "..."
		} else {
			displayText = string(runes)
		}

		res = append(res, fmt.Sprintf("Text: \"%s\". Font: %s. FontSize: %.2fpt. X: %.2fmm. Y: %.2fmm.", displayText, text.Font, text.FontSize, x, y))
	}

	return res, nil
}

func cleanFontName(font string) string {
	if parts := strings.Split(font, "+"); len(parts) > 1 {
		font = parts[1]
	}

	font = strings.ToLower(font)

	font = strings.Split(font, "-")[0]

	return font
}

// FindFontsCoefs считает коэффициент для шрифтов и возрвращает мапу c коэффициентом для каждого шрифта. Эта функция нужна из-за того, что библиотека ledongthuc/pdf показывает координаты левого нижнего угла текста, а typst делает это с левого верхнего угла.
func FindFontsCoefs(page pdf.Page) map[string]float64 {
	coefs := make(map[string]float64)
	resources := page.V.Key("Resources")
	if resources.IsNull() {
		return coefs
	}

	fonts := resources.Key("Font")
	if fonts.IsNull() {
		return coefs
	}

	for _, key := range fonts.Keys() {
		fontVal := fonts.Key(key)

		rawName := fontVal.Key("BaseFont").Name()

		if rawName == "" {
			continue
		}

		cleanName := cleanFontName(rawName)

		if _, exist := coefs[cleanName]; exist {
			continue
		}
		fontDesc := fontVal.Key("FontDescriptor")

		if fontDesc.IsNull() {
			descendants := fontVal.Key("DescendantFonts")
			if !descendants.IsNull() && descendants.Len() > 0 {
				fontDesc = descendants.Index(0).Key("FontDescriptor")
			}
		}
		coef := 0.7

		if !fontDesc.IsNull() {
			capHeight := fontDesc.Key("CapHeight").Float64()
			if capHeight > 0 {
				coef = capHeight / 1000.0
			}
		}

		coefs[cleanName] = coef
	}
	return coefs
}
