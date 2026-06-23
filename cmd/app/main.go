package main

import (
	"fmt"
	"os"

	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/decliner"
	"github.com/syncpath/go-awards/internal/differ"
	"github.com/syncpath/go-awards/internal/finder"
	"github.com/syncpath/go-awards/internal/generator"
	"github.com/syncpath/go-awards/internal/parser"
)

func main() {
	err := decliner.InitDecliner("Rules/rules.json")
	if err != nil {
		fmt.Println("Ошибка инициализации правил Petrovich:", err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig("certificates.json")
	if err != nil {
		fmt.Println("ошибка загрузки конфига")
		os.Exit(1)
	}

	table, err := parser.ParseCSV("NN_reports.csv")
	if err != nil {
		fmt.Println("ошибка загрузки таблицы")
		os.Exit(1)
	}

	result, err := generator.GenerateTypst(cfg, table)
	if err != nil {
		fmt.Println("ошибка генерации typst файла:", err)
		os.Exit(1)
	}

	outputPath := "certificates.typ"
	err = os.WriteFile(outputPath, []byte(result), 0o644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		os.Exit(1)
	}
	fmt.Println("Удачно сгенерировано")

	fields, err := finder.FindFields("Серт.pdf")
	if err != nil {
		fmt.Println("Ошибка поиска полей pdf:", err)
		os.Exit(1)
	}
	for _, field := range fields {
		fmt.Println(field)
	}

	err = differ.DiffPdf("Серт.pdf", "certificates.pdf")
	if err != nil {
		fmt.Println("Ошибка диффера: ", err)
		os.Exit(1)
	}
}
