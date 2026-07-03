package ui

import (
	"fmt"
	"path/filepath"

	"github.com/syncpath/go-awards/internal/config"
)

func (c *UIController) RefreshFiles() {
	// JSON шаблоны из materials/templates
	jsons, _ := filepath.Glob(filepath.Join(templatesDir, "*.json"))
	var jsonNames []string
	for _, p := range jsons {
		jsonNames = append(jsonNames, filepath.Base(p))
	}
	c.jsonsSelect.Options = jsonNames
	c.jsonsSelect.Refresh()

	// CSV таблицы из materials/tables
	csvs, _ := filepath.Glob(filepath.Join(tablesDir, "*.csv"))
	var csvNames []string
	for _, p := range csvs {
		csvNames = append(csvNames, filepath.Base(p))
	}
	c.csvsSelect.Options = csvNames
	c.csvsSelect.Refresh()

	// Diff backgrounds из materials/examples
	diffs, _ := filepath.Glob(filepath.Join(examplesDir, "*.pdf"))
	var diffNames []string
	for _, p := range diffs {
		diffNames = append(diffNames, filepath.Base(p))
	}
	c.diffBacksSelect.Options = diffNames
	c.diffBacksSelect.Refresh()

	// Finder картинки из examples
	c.findBacksSelect.Options = diffNames
	c.findBacksSelect.Refresh()
}

func (c *UIController) LoadTemplateToForm(name string) {
	if name == "" {
		return
	}
	cfg, err := config.LoadConfig(config.RealFS{}, templatesDir, name)
	if err != nil {
		c.output.SetText(fmt.Sprintf("Ошибка загрузки шаблона для редактирования: %v", err))
		return
	}

	c.jsonFileEntry.SetText(name)
	c.titleEntry.SetText(cfg.Title)
	c.widthEntry.SetText(fmt.Sprintf("%.2f", cfg.Width))
	c.heightEntry.SetText(fmt.Sprintf("%.2f", cfg.Height))

	c.bgImgEntry.SetText(cfg.Bg.BackgroundImage)
	if cfg.Bg.BgScale != 0 {
		c.bgScaleEntry.SetText(fmt.Sprintf("%.2f", cfg.Bg.BgScale))
	} else {
		c.bgScaleEntry.SetText("")
	}
	if cfg.Bg.BgDx != 0 {
		c.bgDxEntry.SetText(fmt.Sprintf("%.2f", cfg.Bg.BgDx))
	} else {
		c.bgDxEntry.SetText("")
	}
	if cfg.Bg.BgDy != 0 {
		c.bgDyEntry.SetText(fmt.Sprintf("%.2f", cfg.Bg.BgDy))
	} else {
		c.bgDyEntry.SetText("")
	}

	// очищаем существующие поля UI
	c.fieldsUI = nil
	c.fieldsAccordion.Items = nil

	// загрузка полей
	for _, f := range cfg.Fields {
		c.AddField(f)
	}
	c.fieldsAccordion.Refresh()
}
