package ui

import (
	"path/filepath"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const materialsDir = "materials"

var (
	templatesDir = filepath.Join(materialsDir, "templates")
	tablesDir    = filepath.Join(materialsDir, "tables")
	examplesDir  = filepath.Join(materialsDir, "examples")
	diffDir      = filepath.Join(materialsDir, "diff")
)

type FieldUI struct {
	idEntry         *widget.Entry
	valEntry        *widget.Entry
	xEntry          *widget.Entry
	yEntry          *widget.Entry
	widthEntry      *widget.Entry
	fontEntry       *widget.Entry
	fontTypeEntry   *widget.Entry
	fontSizeEntry   *widget.Entry
	colorEntry      *widget.Entry
	alignEntry      *widget.Entry
	placeAlignEntry *widget.Entry
	leadingEntry    *widget.Entry
	spacingEntry    *widget.Entry
	trackingEntry   *widget.Entry
	indentEntry     *widget.Entry
	indentValEntry  *widget.Entry
	accordionItem   *widget.AccordionItem
}

type UIController struct {
	output          *ReadOnlyEntry
	finderRes       *ReadOnlyEntry
	jsonsSelect     *widget.Select
	csvsSelect      *widget.Select
	diffBacksSelect *widget.Select
	findBacksSelect *widget.Select

	imageViewer *canvas.Image
	imageLayout *ZoomLayout
	imageScroll *container.Scroll

	genCompBut *widget.Button
	diffBut    *widget.Button
	findBut    *widget.Button
	refreshBut *widget.Button

	// zoom
	zoomLevel  float32
	zoomInBut  *widget.Button
	zoomOutBut *widget.Button

	// template creator
	jsonFileEntry *widget.Entry
	titleEntry    *widget.Entry
	widthEntry    *widget.Entry
	heightEntry   *widget.Entry
	bgImgEntry    *widget.Entry
	bgScaleEntry  *widget.Entry
	bgDxEntry     *widget.Entry
	bgDyEntry     *widget.Entry

	// удобный список с выпадающими полями для полей шаблона
	fieldsAccordion *widget.Accordion
	fieldsUI        []*FieldUI
}

func (c *UIController) setUIEnabled(enabled bool) {
	if enabled {
		c.genCompBut.Enable()
		c.diffBut.Enable()
		c.findBut.Enable()
		c.refreshBut.Enable()
	} else {
		c.genCompBut.Disable()
		c.diffBut.Disable()
		c.findBut.Disable()
		c.refreshBut.Disable()
	}
}
