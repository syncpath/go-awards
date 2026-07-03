// Package ui пакет с пользовательским интерфейсом
package ui

import (
	"fmt"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/syncpath/go-awards/internal/approot"
	"github.com/syncpath/go-awards/internal/config"
	"github.com/syncpath/go-awards/internal/decliner"
)

func Run() {
	// переходим в корень проекта относительно бинарника
	_, rootErr := approot.Setup()

	a := app.New()
	w := a.NewWindow("go-awards app")
	w.SetFullScreen(true)

	//===============
	// Правая сторона
	//===============

	// finder (ReadOnlyEntry для сохранения яркого цвета текста)
	finderRes := NewReadOnlyEntry()

	scrollFinder := container.NewVScroll(finderRes)
	scrollFinder.SetMinSize(fyne.NewSize(0, 188))

	finderWidget := widget.NewCard("", "Finder result", scrollFinder)

	// json список выпадающий
	jsons := widget.NewSelect([]string{}, func(chosen string) {})
	jsonWidget := widget.NewCard("", "JSON templates", jsons)

	// csv таблицы выпадающий список
	csvs := widget.NewSelect([]string{}, func(chosen string) {})
	csvWidget := widget.NewCard("", "CSV tables", csvs)

	// фоны для diff выпадающий список
	diffBacks := widget.NewSelect([]string{}, func(chosen string) {})
	diffWidget := widget.NewCard("", "Diff backgrounds", diffBacks)

	// pdf для finder выпадающий список
	findBacks := widget.NewSelect([]string{}, func(chosen string) {})
	findWidget := widget.NewCard("", "Finder images", findBacks)

	// добавляем все в бокс верхний
	topBox := container.NewGridWithColumns(2, jsonWidget, csvWidget, diffWidget, findWidget)

	// cоздание шаблона в центре правой части
	jsonFileEntry := widget.NewEntry()
	jsonFileEntry.SetPlaceHolder("new_template.json")
	titleEntry := widget.NewEntry()
	widthEntry := widget.NewEntry()
	heightEntry := widget.NewEntry()

	bgImgEntry := widget.NewEntry()
	bgScaleEntry := widget.NewEntry()
	bgDxEntry := widget.NewEntry()
	bgDyEntry := widget.NewEntry()

	templateForm := widget.NewForm(
		widget.NewFormItem("File Name *", jsonFileEntry),
		widget.NewFormItem("Title *", titleEntry),
		widget.NewFormItem("Width (mm) *", widthEntry),
		widget.NewFormItem("Height (mm) *", heightEntry),
		widget.NewFormItem("Bg Image", bgImgEntry),
		widget.NewFormItem("Bg Scale (%)", bgScaleEntry),
		widget.NewFormItem("Bg DX (mm)", bgDxEntry),
		widget.NewFormItem("Bg DY (mm)", bgDyEntry),
	)

	fieldsAccordion := widget.NewAccordion()
	addFieldBut := widget.NewButton("Add Field", func() {})
	saveTemplateBut := widget.NewButton("Save Template", func() {})

	creatorContainer := container.NewVBox(
		widget.NewLabel("General Settings"),
		templateForm,
		widget.NewLabel("Fields Configuration"),
		fieldsAccordion,
		addFieldBut,
		widget.NewSeparator(),
		saveTemplateBut,
	)

	creatorScroll := container.NewVScroll(creatorContainer)
	creatorScroll.SetMinSize(fyne.NewSize(0, 250))
	creatorWidget := widget.NewCard("", "Template Creator", creatorScroll)

	rightArea := container.NewBorder(topBox, finderWidget, nil, nil, creatorWidget)

	//#############
	//Левая сторона
	//#############

	// блок с кнопками
	genCompBut := widget.NewButton("Generate and Compile", func() {})
	diffBut := widget.NewButton("Generate and Diff", func() {})
	findBut := widget.NewButton("Start Finder", func() {})
	refreshBut := widget.NewButton("Refresh", func() {})

	buttonBox := container.NewGridWithColumns(4, genCompBut, diffBut, findBut, refreshBut)

	// output
	output := NewReadOnlyEntry()

	scrollOutput := container.NewVScroll(output)
	scrollOutput.SetMinSize(fyne.NewSize(0, 145))

	outputWidget := widget.NewCard("", "Output", scrollOutput)

	bottomBox := container.NewBorder(buttonBox, nil, nil, nil, outputWidget)

	// создаем начальную картинку, белый квадрат
	rect := image.Rect(0, 0, 1000, 1000)
	blankImg := image.NewRGBA(rect)
	for x := range 1000 {
		for y := range 1000 {
			blankImg.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}

	imgViewer := canvas.NewImageFromImage(blankImg)
	imgViewer.FillMode = canvas.ImageFillContain

	// контейнер с кастомным лейаутом для переопределения MinSize картинки
	imageLayout := &ZoomLayout{size: fyne.NewSize(1000, 1000)}
	imgWrapper := container.New(imageLayout, imgViewer)

	imgScroll := container.NewScroll(imgWrapper)
	imgScroll.SetMinSize(fyne.NewSize(0, 300))
	imgCard := widget.NewCard("", "Preview", imgScroll)

	// кнопки зума
	zoomIn := widget.NewButton("+", func() {})
	zoomOut := widget.NewButton("-", func() {})
	zoomBar := container.NewHBox(layout.NewSpacer(), zoomOut, zoomIn)

	leftArea := container.NewBorder(zoomBar, bottomBox, nil, nil, imgCard)

	split := container.NewHSplit(leftArea, rightArea)
	split.Offset = 0.6

	// cоздаем и связываем контроллер
	ctrl := &UIController{
		output:          output,
		finderRes:       finderRes,
		jsonsSelect:     jsons,
		csvsSelect:      csvs,
		diffBacksSelect: diffBacks,
		findBacksSelect: findBacks,
		imageViewer:     imgViewer,
		imageLayout:     imageLayout,
		imageScroll:     imgScroll,
		genCompBut:      genCompBut,
		diffBut:         diffBut,
		findBut:         findBut,
		refreshBut:      refreshBut,
		zoomLevel:       1.0,
		zoomInBut:       zoomIn,
		zoomOutBut:      zoomOut,
		jsonFileEntry:   jsonFileEntry,
		titleEntry:      titleEntry,
		widthEntry:      widthEntry,
		heightEntry:     heightEntry,
		bgImgEntry:      bgImgEntry,
		bgScaleEntry:    bgScaleEntry,
		bgDxEntry:       bgDxEntry,
		bgDyEntry:       bgDyEntry,
		fieldsAccordion: fieldsAccordion,
	}

	genCompBut.OnTapped = ctrl.GenerateAndCompile
	diffBut.OnTapped = ctrl.Diff
	findBut.OnTapped = ctrl.StartFinder
	refreshBut.OnTapped = ctrl.RefreshFiles

	zoomIn.OnTapped = ctrl.ZoomIn
	zoomOut.OnTapped = ctrl.ZoomOut

	saveTemplateBut.OnTapped = ctrl.SaveTemplate
	addFieldBut.OnTapped = func() {
		ctrl.AddField(config.Field{})
	}

	// при выборе шаблона загружаем его настройки
	jsons.OnChanged = func(chosen string) {
		ctrl.LoadTemplateToForm(chosen)
	}

	// загружаем список файлов при старте
	ctrl.RefreshFiles()

	if rootErr != nil {
		output.SetText(fmt.Sprintf("не удалось определить корень проекта: %v", rootErr))
	}

	// инициализируем Petrovich
	err := decliner.InitDecliner("Rules/rules.json")
	if err != nil {
		output.SetText(fmt.Sprintf("ошибка инициализации правил склонения Petrovich: %v", err))
	} else if rootErr == nil {
		output.SetText("приложение успешно запущено.")
	}

	w.SetContent(split)
	w.ShowAndRun()
}
