package main

import (
	"fmt"
	"os"

	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/finder"
	"github.com/syncpath/go-awards/internal/generator"
	"github.com/syncpath/go-awards/internal/parser"
)

func main() {
	cfg, err := config.LoadConfig("certificates.json")
	if err != nil {
		fmt.Println("ошибка загрузки конфига")
		os.Exit(1)
	}

	table, err := parser.ParseCSV("example.csv")
	if err != nil {
		fmt.Println("ошибка загрузки таблицы")
		os.Exit(1)
	}

	result, err := generator.GenerateTypst(cfg, table)
	if err != nil {
		fmt.Println("ошибка генерации typst файла")
		os.Exit(1)
	}

	outputPath := "certificates.typ"
	err = os.WriteFile(outputPath, []byte(result), 0o644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		os.Exit(1)
	}
	fmt.Println("Удачно сгенерировано")

	fields, err := finder.FindFields("certificates.pdf")
	if err != nil {
		fmt.Println("Ошибка поиска полей pdf:", err)
		os.Exit(1)
	}
	for _, field := range fields {
		fmt.Println(field)
	}
}
