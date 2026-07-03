package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/syncpath/go-awards/internal/config"
)

func (c *UIController) SaveTemplate() {
	fileName := strings.TrimSpace(c.jsonFileEntry.Text)
	title := strings.TrimSpace(c.titleEntry.Text)
	widthStr := strings.TrimSpace(c.widthEntry.Text)
	heightStr := strings.TrimSpace(c.heightEntry.Text)

	if fileName == "" || title == "" || widthStr == "" || heightStr == "" {
		c.output.SetText("ошибка сохранения: заполните все обязательные поля шаблона (*)")
		return
	}

	if !strings.HasSuffix(fileName, ".json") {
		fileName += ".json"
	}

	var width, height float64
	_, errW := fmt.Sscanf(widthStr, "%f", &width)
	_, errH := fmt.Sscanf(heightStr, "%f", &height)
	if errW != nil || errH != nil {
		c.output.SetText("ошибка сохранения: ширина и высота шаблона должны быть числами")
		return
	}

	// опциональные поля конфига
	bgImg := strings.TrimSpace(c.bgImgEntry.Text)
	var bgScale, bgDx, bgDy float64
	if strings.TrimSpace(c.bgScaleEntry.Text) != "" {
		_, _ = fmt.Sscanf(c.bgScaleEntry.Text, "%f", &bgScale)
	}
	if strings.TrimSpace(c.bgDxEntry.Text) != "" {
		_, _ = fmt.Sscanf(c.bgDxEntry.Text, "%f", &bgDx)
	}
	if strings.TrimSpace(c.bgDyEntry.Text) != "" {
		_, _ = fmt.Sscanf(c.bgDyEntry.Text, "%f", &bgDy)
	}

	// прочитать поля из fieldsUI
	var fields []config.Field
	for _, fui := range c.fieldsUI {
		id := strings.TrimSpace(fui.idEntry.Text)
		val := strings.TrimSpace(fui.valEntry.Text)
		xStr := strings.TrimSpace(fui.xEntry.Text)
		yStr := strings.TrimSpace(fui.yEntry.Text)
		wStr := strings.TrimSpace(fui.widthEntry.Text)
		font := strings.TrimSpace(fui.fontEntry.Text)
		fontSizeStr := strings.TrimSpace(fui.fontSizeEntry.Text)
		color := strings.TrimSpace(fui.colorEntry.Text)

		// валидация только обязательных полей
		if id == "" || val == "" || wStr == "" || font == "" || fontSizeStr == "" {
			c.output.SetText(fmt.Sprintf("ошибка сохранения: заполните все обязательные поля (*) для элемента '%s'", id))
			return
		}

		var x, y, wVal, fontSize float64
		if xStr != "" {
			_, errX := fmt.Sscanf(xStr, "%f", &x)
			if errX != nil {
				c.output.SetText(fmt.Sprintf("ошибка сохранения: X для элемента '%s' должно быть числом", id))
				return
			}
		}
		if yStr != "" {
			_, errY := fmt.Sscanf(yStr, "%f", &y)
			if errY != nil {
				c.output.SetText(fmt.Sprintf("ошибка сохранения: Y для элемента '%s' должно быть числом", id))
				return
			}
		}
		_, errWVal := fmt.Sscanf(wStr, "%f", &wVal)
		_, errFS := fmt.Sscanf(fontSizeStr, "%f", &fontSize)
		if errWVal != nil || errFS != nil {
			c.output.SetText(fmt.Sprintf("ошибка сохранения: числовые параметры (Width, Font Size) элемента '%s' содержат некорректные символы", id))
			return
		}

		// чтение 8 опциональный полей из Field
		fontType := strings.TrimSpace(fui.fontTypeEntry.Text)
		if fontType == "" {
			fontType = "regular" // дефолтное значение для Typst
		}
		align := strings.TrimSpace(fui.alignEntry.Text)
		placeAlign := strings.TrimSpace(fui.placeAlignEntry.Text)
		indent := strings.TrimSpace(fui.indentEntry.Text)

		var leading, spacing, tracking, indentVal float64
		if strings.TrimSpace(fui.leadingEntry.Text) != "" {
			_, _ = fmt.Sscanf(fui.leadingEntry.Text, "%f", &leading)
		}
		if strings.TrimSpace(fui.spacingEntry.Text) != "" {
			_, _ = fmt.Sscanf(fui.spacingEntry.Text, "%f", &spacing)
		}
		if strings.TrimSpace(fui.trackingEntry.Text) != "" {
			_, _ = fmt.Sscanf(fui.trackingEntry.Text, "%f", &tracking)
		}
		if strings.TrimSpace(fui.indentValEntry.Text) != "" {
			_, _ = fmt.Sscanf(fui.indentValEntry.Text, "%f", &indentVal)
		}

		fields = append(fields, config.Field{
			ID:          id,
			Value:       val,
			X:           x,
			Y:           y,
			Width:       wVal,
			Font:        font,
			FontType:    fontType,
			FontSize:    fontSize,
			Color:       color,
			Align:       align,
			PlaceAlign:  placeAlign,
			Leading:     leading,
			Spacing:     spacing,
			Tracking:    tracking,
			Indent:      indent,
			IndentValue: indentVal,
		})
	}

	cfg := config.TemplateConfig{
		Title:  title,
		Width:  width,
		Height: height,
		Bg: config.Background{
			BackgroundImage: bgImg,
			BgScale:         bgScale,
			BgDx:            bgDx,
			BgDy:            bgDy,
		},
		Fields: fields,
	}

	importJSON, errMarshal := json.MarshalIndent(cfg, "", "  ")
	if errMarshal != nil {
		c.output.SetText(fmt.Sprintf("ошибка сохранения: %v", errMarshal))
		return
	}

	outputPath := filepath.Join(templatesDir, fileName)
	errWrite := os.WriteFile(outputPath, importJSON, 0o644)
	if errWrite != nil {
		c.output.SetText(fmt.Sprintf("ошибка записи файла %s: %v", outputPath, errWrite))
		return
	}

	c.output.SetText(fmt.Sprintf("Шаблон %s успешно сохранен с %d полями!", fileName, len(fields)))

	// обновить список файлов
	c.RefreshFiles()
}
