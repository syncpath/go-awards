//go:build integration

package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/syncpath/go-awards/internal/approot"
	"github.com/syncpath/go-awards/internal/config"
)

func TestMain(m *testing.M) {
	if _, err := approot.Setup(); err != nil {
		panic("не удалось определить корень проекта для теста: " + err.Error())
	}
	test.NewApp()
	os.Exit(m.Run())
}

func newTestController() *UIController {
	return &UIController{
		output:          NewReadOnlyEntry(),
		finderRes:       NewReadOnlyEntry(),
		jsonsSelect:     widget.NewSelect(nil, func(string) {}),
		csvsSelect:      widget.NewSelect(nil, func(string) {}),
		diffBacksSelect: widget.NewSelect(nil, func(string) {}),
		findBacksSelect: widget.NewSelect(nil, func(string) {}),
		jsonFileEntry:   widget.NewEntry(),
		titleEntry:      widget.NewEntry(),
		widthEntry:      widget.NewEntry(),
		heightEntry:     widget.NewEntry(),
		bgImgEntry:      widget.NewEntry(),
		bgScaleEntry:    widget.NewEntry(),
		bgDxEntry:       widget.NewEntry(),
		bgDyEntry:       widget.NewEntry(),
		fieldsAccordion: widget.NewAccordion(),
	}
}

// полный цикл создания шаблона через UI
func TestSaveAndLoadTemplateRoundTrip(t *testing.T) {
	c := newTestController()

	const fileName = "ui_roundtrip_test.json"
	outputPath := filepath.Join(templatesDir, fileName)
	t.Cleanup(func() { os.Remove(outputPath) })

	// заполняем общие настройки шаблона
	c.jsonFileEntry.SetText(fileName)
	c.titleEntry.SetText("Тестовый шаблон")
	c.widthEntry.SetText("210.0")
	c.heightEntry.SetText("297.0")
	c.bgImgEntry.SetText("images/bg.png")
	c.bgScaleEntry.SetText("105.0")

	// добавляем одно поле со всеми ключевыми параметрами
	c.AddField(config.Field{
		ID:       "fio",
		Value:    "@name(\"ФИО\", дательный, false, 1)",
		Width:    190.0,
		Font:     "PT Sans",
		FontType: "bold",
		FontSize: 32.0,
		Leading:  5.64,
		X:        10.0,
		Y:        20.0,
	})

	c.SaveTemplate()

	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("файл шаблона не был создан: %v", err)
	}
	if strings.HasPrefix(c.output.Text, "ошибка") {
		t.Fatalf("SaveTemplate вернул ошибку: %s", c.output.Text)
	}

	// загружаем сохранённый шаблон обратно в форму
	c2 := newTestController()
	c2.LoadTemplateToForm(fileName)

	if c2.titleEntry.Text != "Тестовый шаблон" {
		t.Errorf("title: получил %q", c2.titleEntry.Text)
	}
	if c2.widthEntry.Text != "210.00" {
		t.Errorf("width: получил %q", c2.widthEntry.Text)
	}
	if len(c2.fieldsUI) != 1 {
		t.Fatalf("ожидал 1 поле после загрузки, получил %d", len(c2.fieldsUI))
	}
	f := c2.fieldsUI[0]
	if f.idEntry.Text != "fio" {
		t.Errorf("field id: получил %q", f.idEntry.Text)
	}
	if f.fontSizeEntry.Text != "32.00" {
		t.Errorf("field font size: получил %q", f.fontSizeEntry.Text)
	}
	if f.leadingEntry.Text != "5.64" {
		t.Errorf("field leading: получил %q", f.leadingEntry.Text)
	}
}

// проверяет, при незаполненных обязательных полях ошибка
func TestSaveTemplateValidation(t *testing.T) {
	c := newTestController()

	// не заполняем ничего
	c.SaveTemplate()

	if !strings.HasPrefix(c.output.Text, "ошибка") {
		t.Fatalf("ожидал ошибку, получил: %q", c.output.Text)
	}
	if _, err := os.Stat(filepath.Join(templatesDir, ".json")); err == nil {
		t.Fatalf("файл не должен был быть создан при ошибке валидации")
	}
}

// проверяет отказ при нечисловой ширине/высоте
func TestSaveTemplateBadNumber(t *testing.T) {
	c := newTestController()
	c.jsonFileEntry.SetText("bad_number_test")
	c.titleEntry.SetText("x")
	c.widthEntry.SetText("не число")
	c.heightEntry.SetText("297")

	c.SaveTemplate()

	if !strings.Contains(c.output.Text, "числами") {
		t.Fatalf("ожидал ошибку о числах, получил: %q", c.output.Text)
	}
	t.Cleanup(func() { os.Remove(filepath.Join(templatesDir, "bad_number_test.json")) })
	if _, err := os.Stat(filepath.Join(templatesDir, "bad_number_test.json")); err == nil {
		t.Fatalf("файл не должен создаваться при некорректных числах")
	}
}

// проверяет, что списки файлов подтягиваются из папок проекта
func TestRefreshFiles(t *testing.T) {
	c := newTestController()
	c.RefreshFiles()

	if len(c.jsonsSelect.Options) == 0 {
		t.Error("ожидал непустой список JSON-шаблонов")
	}
	// Finder и Diff используют один и тот же список PDF из examples
	if len(c.findBacksSelect.Options) != len(c.diffBacksSelect.Options) {
		t.Errorf("списки finder и diff должны совпадать: %d vs %d",
			len(c.findBacksSelect.Options), len(c.diffBacksSelect.Options))
	}
}

// проверяет, что заголовок аккордеона обновляется по ID
func TestAddFieldRenamesHeader(t *testing.T) {
	c := newTestController()
	c.AddField(config.Field{})

	if len(c.fieldsUI) != 1 {
		t.Fatalf("ожидал 1 поле, получил %d", len(c.fieldsUI))
	}
	item := c.fieldsUI[0].accordionItem
	if item.Title != "new_field" {
		t.Errorf("заголовок пустого поля: получил %q, ожидал new_field", item.Title)
	}

	c.fieldsUI[0].idEntry.SetText("myid")
	if item.Title != "myid" {
		t.Errorf("заголовок после ввода ID: получил %q, ожидал myid", item.Title)
	}
}
