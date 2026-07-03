package main

import (
	"fmt"
	"os"

	"github.com/syncpath/go-awards/internal/approot"
	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/decliner"
	"github.com/syncpath/go-awards/internal/differ"
	"github.com/syncpath/go-awards/internal/finder"
	"github.com/syncpath/go-awards/internal/generator"
	"github.com/syncpath/go-awards/internal/parser"
)

func main() {
	if _, err := approot.Setup(); err != nil {
		fmt.Println("Ошибка определения корня проекта:", err)
		os.Exit(1)
	}

	err := decliner.InitDecliner("Rules/rules.json")
	if err != nil {
		fmt.Println("Ошибка инициализации правил Petrovich:", err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(config.RealFS{}, "materials/templates", "certificates.json")
	if err != nil {
		fmt.Println("ошибка загрузки конфига")
		os.Exit(1)
	}

	table, err := parser.ParseCSV("materials/tables", "NN_reports.csv")
	if err != nil {
		fmt.Println("ошибка загрузки таблицы")
		os.Exit(1)
	}

	result, err := generator.GenerateTypst(cfg, table)
	if err != nil {
		fmt.Println("ошибка генерации typst файла:", err)
		os.Exit(1)
	}

	outputPath := "materials/certificates.typ"
	err = os.WriteFile(outputPath, []byte(result), 0o644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		os.Exit(1)
	}
	fmt.Println("Удачно сгенерировано")

	fields, err := finder.ReadPdf("materials/examples", "blagG.pdf")
	if err != nil {
		fmt.Println("Ошибка поиска полей pdf:", err)
		os.Exit(1)
	}
	fmt.Printf("Height %v, Width: %v\n", fields.Height, fields.Width)
	for _, field := range fields.Texts {
		fmt.Printf("Teкст: %s. Шрифт: %s. Размер шрифта: %v pt. X: %v mm. Y: %v mm. Отступ: %v mm\n", field.Text, field.Font, field.FontSize, field.X, field.Y, field.Leading)
	}

	err = differ.DiffPdf("Серт.pdf", "materials/certificates.pdf")
	if err != nil {
		fmt.Println("Ошибка диффера: ", err)
		os.Exit(1)
	}
}
