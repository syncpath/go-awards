package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/differ"
	"github.com/syncpath/go-awards/internal/finder"
	"github.com/syncpath/go-awards/internal/generator"
	"github.com/syncpath/go-awards/internal/parser"
)

func (c *UIController) GenerateAndCompile() {
	jsonSel := c.jsonsSelect.Selected
	csvSel := c.csvsSelect.Selected

	if jsonSel == "" || csvSel == "" {
		c.output.SetText("ошибка генерации typst файла: не выбран json шаблон или csv таблица")
		return
	}

	c.setUIEnabled(false)
	c.output.SetText("генерация и компиляция...")

	go func() {
		defer func() {
			fyne.Do(func() {
				c.setUIEnabled(true)
			})
		}()

		// сгенерировать typst
		cfg, err := config.LoadConfig(config.RealFS{}, templatesDir, jsonSel)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка загрузки шаблона: %v", err)) })
			return
		}

		table, err := parser.ParseCSV(tablesDir, csvSel)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка загрузки таблицы: %v", err)) })
			return
		}

		result, err := generator.GenerateTypst(cfg, table)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка генерации typst: %v", err)) })
			return
		}

		baseName := strings.TrimSuffix(jsonSel, filepath.Ext(jsonSel))
		typPath := filepath.Join(materialsDir, baseName+".typ")
		pdfPath := filepath.Join(materialsDir, baseName+".pdf")

		// убедиться, что diff директория существует
		_ = os.MkdirAll(diffDir, 0o755)
		pngPath := filepath.Join(diffDir, baseName+".png")

		err = os.WriteFile(typPath, []byte(result), 0o644)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка записи файла %s: %v", typPath, err)) })
			return
		}

		// скомпилить в pdf
		cmdPdf := exec.Command("typst", "compile", typPath, pdfPath)
		outputPdf, err := cmdPdf.CombinedOutput()
		if err != nil {
			fyne.Do(func() {
				c.output.SetText(fmt.Sprintf("ошибка компиляции PDF:\n%s\n%v", string(outputPdf), err))
			})
			return
		}

		// скомпилить в PNG
		cmdPng := exec.Command("typst", "compile", "--pages", "1", typPath, pngPath)
		outputPng, err := cmdPng.CombinedOutput()
		if err != nil {
			fyne.Do(func() {
				c.output.SetText(fmt.Sprintf("ошибка компиляции PNG:\n%s\n%v", string(outputPng), err))
			})
			return
		}

		// обновление ресурса для кэша fyne
		res, loadErr := fyne.LoadResourceFromPath(pngPath)

		// обновить превью в UI
		fyne.Do(func() {
			if loadErr == nil {
				c.imageViewer.Resource = res
				c.imageViewer.File = ""
				c.imageViewer.Refresh()
				c.output.SetText(fmt.Sprintf("скомпилировано успешно!\nPDF сохранен как: %s\nпревью PNG сохранено как: %s", pdfPath, pngPath))
			} else {
				c.output.SetText(fmt.Sprintf("скомпилировано успешно, но не удалось загрузить превью: %v", loadErr))
			}
		})
	}()
}

func (c *UIController) Diff() {
	jsonSel := c.jsonsSelect.Selected
	csvSel := c.csvsSelect.Selected
	diffBackSel := c.diffBacksSelect.Selected

	if jsonSel == "" || csvSel == "" || diffBackSel == "" {
		c.output.SetText("ошибка диффа: проверьте, что выбраны JSON шаблон, CSV таблица и PDF для диффа")
		return
	}

	c.setUIEnabled(false)
	c.output.SetText("генерация, компиляция и сравнение...")

	go func() {
		defer func() {
			fyne.Do(func() {
				c.setUIEnabled(true)
			})
		}()

		// генерация и компиляция PDF
		cfg, err := config.LoadConfig(config.RealFS{}, templatesDir, jsonSel)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка загрузки шаблона: %v", err)) })
			return
		}

		table, err := parser.ParseCSV(tablesDir, csvSel)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка загрузки таблицы: %v", err)) })
			return
		}

		result, err := generator.GenerateTypst(cfg, table)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка генерации typst: %v", err)) })
			return
		}

		baseName := strings.TrimSuffix(jsonSel, filepath.Ext(jsonSel))
		typPath := filepath.Join(materialsDir, baseName+".typ")
		pdfPath := filepath.Join(materialsDir, baseName+".pdf")

		err = os.WriteFile(typPath, []byte(result), 0o644)
		if err != nil {
			fyne.Do(func() { c.output.SetText(fmt.Sprintf("ошибка записи файла %s: %v", typPath, err)) })
			return
		}

		// компиляция PDF
		cmdPdf := exec.Command("typst", "compile", typPath, pdfPath)
		outputPdf, err := cmdPdf.CombinedOutput()
		if err != nil {
			fyne.Do(func() {
				c.output.SetText(fmt.Sprintf("Ошибка компиляции PDF:\n%s\n%v", string(outputPdf), err))
			})
			return
		}

		// вызов differ.DiffPdf
		err = differ.DiffPdf(diffBackSel, pdfPath)
		if err != nil {
			fyne.Do(func() {
				c.output.SetText(fmt.Sprintf("Ошибка сравнения (diff-pdf): %v", err))
			})
			return
		}

		// выходные данные дифера materials/diff/output-diff-1.png
		diffPngPath := filepath.Join(diffDir, "output-diff-1.png")

		// обновление ресурса для кэша fyne
		res, loadErr := fyne.LoadResourceFromPath(diffPngPath)

		// обновить превью в UI
		fyne.Do(func() {
			if loadErr == nil {
				c.imageViewer.Resource = res
				c.imageViewer.File = ""
				c.imageViewer.Refresh()
				c.output.SetText("дифф выполнен успешно!\nИзображение сравнения загружено.")
			} else {
				c.output.SetText(fmt.Sprintf("дифф выполнен успешно, но не удалось загрузить превью: %v", loadErr))
			}
		})
	}()
}

func (c *UIController) StartFinder() {
	findBackSel := c.findBacksSelect.Selected

	if findBackSel == "" {
		c.finderRes.SetText("ошибка: не выбран файл для поиска (Finder image)")
		return
	}

	c.setUIEnabled(false)
	c.finderRes.SetText("сканирование макета...")

	go func() {
		defer func() {
			fyne.Do(func() {
				c.setUIEnabled(true)
			})
		}()

		fields, err := finder.ReadPdf(examplesDir, findBackSel)
		if err != nil {
			fyne.Do(func() {
				c.finderRes.SetText(fmt.Sprintf("ошибка поиска полей pdf: %v", err))
			})
			return
		}

		// вывод файндера
		var sb strings.Builder
		fmt.Fprintf(&sb, "Height: %v mm, Width: %v mm\n\n", fields.Height, fields.Width)
		for i, field := range fields.Texts {
			fmt.Fprintf(&sb, "[%d] Текст: %s\n    Шрифт: %s\n    Размер: %v pt\n    Координаты: X=%v mm, Y=%v mm\n    Межстрочный: %v mm\n\n",
				i+1, field.Text, field.Font, field.FontSize, field.X, field.Y, field.Leading)
		}

		fyne.Do(func() {
			c.finderRes.SetText(sb.String())
		})
	}()
}
