// Package finder Изучает наш макет (например, из брендбука) и выводит всю информацию о нем
package finder

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

const coefmm float64 = 25.4 / 72.0

const defCoef float64 = 0.7

// MergedText вспомогательная структура, нужна, так как ledongthuc/pdf хранит каждый символ как отдельный текст
type MergedText struct {
	Font     string
	FontSize float64
	X        float64
	Y        float64
	S        string
	LastX    float64
}

// PageData вытаскиваем информацию о странице в extractPageData
type PageData struct {
	Texts    []Text
	XMin     float64
	XMax     float64
	YMin     float64
	YMax     float64
	CoefsMap map[string]float64
}

// Text пародия на pdf.Text для PageData
type Text struct {
	Font     string
	FontSize float64
	X, Y     float64
	S        string
}

// TextField информация о строке для итогового результата
type TextField struct {
	Font     string
	FontSize float64
	X, Y     float64
	Text     string
	Leading  float64
}

// PageInfo итоговый результат
type PageInfo struct {
	Height float64
	Width  float64
	Texts  []TextField
}

func ReadPdf(path string, name string) (PageInfo, error) {
	fullpath := filepath.Join(path, name)
	file, reader, err := pdf.Open(fullpath)
	if err != nil {
		return PageInfo{}, fmt.Errorf("ошибка чтения макета: %w", err)
	}

	defer file.Close()

	// Подразумевается, что макет имеет только 1 страницу
	p := reader.Page(1)
	if p.V.IsNull() {
		return PageInfo{}, errors.New("страница пустая")
	}

	data, err := extractPageData(p)
	if err != nil {
		return PageInfo{}, err
	}

	return FindFields(data), nil
}

func extractPageData(p pdf.Page) (PageData, error) {
	textsRaw := p.Content().Text
	var texts []Text
	for _, t := range textsRaw {
		text := Text{S: t.S, X: t.X, Y: t.Y, Font: t.Font, FontSize: t.FontSize}
		texts = append(texts, text)
	}

	// Параметры страницы
	mediaBox := p.V.Key("MediaBox")
	xMin, xMax := mediaBox.Index(0).Float64(), mediaBox.Index(2).Float64()
	yMin, yMax := mediaBox.Index(1).Float64(), mediaBox.Index(3).Float64()

	// Коэффициент шрифта
	capHeigths := make(map[string]float64)

	resources := p.V.Key("Resources")
	if resources.IsNull() {
		return PageData{}, errors.New("не найдены ресурсы на странице")
	}

	pfonts := resources.Key("Font")
	if pfonts.IsNull() {
		return PageData{}, errors.New("не найдены шрифты на странице")
	}

	for _, key := range pfonts.Keys() {
		fontVal := pfonts.Key(key)

		rawName := fontVal.Key("BaseFont").Name()

		if rawName == "" {
			continue
		}

		cleanName := CleanFontName(rawName)

		if _, exist := capHeigths[cleanName]; exist {
			continue
		}
		fontDesc := fontVal.Key("FontDescriptor")

		if fontDesc.IsNull() {
			descendants := fontVal.Key("DescendantFonts")
			if !descendants.IsNull() && descendants.Len() > 0 {
				fontDesc = descendants.Index(0).Key("FontDescriptor")
			}
		}
		if !fontDesc.IsNull() {
			capHeigths[cleanName] = fontDesc.Key("CapHeight").Float64()
		}
	}
	return PageData{Texts: texts, XMin: xMin, XMax: xMax, YMin: yMin, YMax: yMax, CoefsMap: FindFontsCoefs(capHeigths)}, nil
}

// FindFields показывает координаты текстовых полей, которые она смогла найти
func FindFields(data PageData) PageInfo {
	var res PageInfo

	var merged []MergedText

	for _, t := range data.Texts {
		if t.S == "" {
			continue
		}

		if len(merged) > 0 {
			last := &merged[len(merged)-1]

			// Шрифты равны, Y равны, расстояние между соседними буквами не больше 20 мм (сравниваем с LastX последней буквы)
			if last.Font == t.Font && last.FontSize == t.FontSize && last.Y == t.Y && math.Abs(t.X-last.LastX) < 20*(1/coefmm) {
				last.S += t.S
				last.LastX = t.X
				continue
			}
		}

		merged = append(merged, MergedText{
			Font:     t.Font,
			FontSize: t.FontSize,
			X:        t.X,
			Y:        t.Y,
			S:        t.S,
			LastX:    t.X,
		})
	}

	height := (data.YMax - data.YMin) * coefmm
	width := (data.XMax - data.XMin) * coefmm

	res.Height, res.Width = round2(height), round2(width)

	for i, text := range merged {
		cleanName := CleanFontName(text.Font)
		coef, ok := data.CoefsMap[cleanName]
		if !ok {
			coef = defCoef
		}

		x := (text.X - data.XMin) * coefmm
		y := (data.YMax - (text.Y + text.FontSize*coef)) * coefmm

		runes := []rune(text.S)
		var displayText string
		if len(runes) > 15 {
			displayText = string(runes[:15]) + "..."
		} else {
			displayText = string(runes)
		}

		var leadingMm float64
		if i > 0 {
			prev := merged[i-1]
			if prev.Y > text.Y {
				distPts := prev.Y - text.Y

				gapPts := distPts - (text.FontSize * coef)
				leadingMm = gapPts * coefmm
			}
		}
		res.Texts = append(res.Texts, TextField{X: round2(x), Y: round2(y), Font: text.Font, FontSize: round2(text.FontSize), Text: displayText, Leading: round2(leadingMm)})
	}

	return res
}

// FindFontsCoefs считает коэффициент для шрифтов и возрвращает мапу c коэффициентом для каждого шрифта. Эта функция нужна из-за того, что библиотека ledongthuc/pdf показывает координаты левого нижнего угла текста, а typst делает это с левого верхнего угла.
func FindFontsCoefs(fontsCapHeigths map[string]float64) map[string]float64 {
	coefs := make(map[string]float64)

	for font := range fontsCapHeigths {
		coef := defCoef

		if fontsCapHeigths[font] > 0 {
			if fontsCapHeigths[font] > 1000.0 {
				coef = fontsCapHeigths[font] / 2048.0
			} else {
				coef = fontsCapHeigths[font] / 1000.0
			}
		}

		if coef < 0.6 || coef > 0.8 {
			coef = defCoef
		}

		coefs[font] = round2(coef)
	}
	return coefs
}

func round2(num float64) float64 {
	return math.Round(num*100) / 100
}

func CleanFontName(font string) string {
	if i := strings.Index(font, "+"); i >= 0 {
		font = font[i+1:]
	}

	return font
}
