package ui

import (
	"fmt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/syncpath/go-awards/internal/config"
)

func (c *UIController) AddField(f config.Field) {
	// обязательные поля
	idEntry := widget.NewEntry()
	idEntry.SetText(f.ID)
	idEntry.SetPlaceHolder("fio")

	valEntry := widget.NewEntry()
	valEntry.SetText(f.Value)
	valEntry.SetPlaceHolder("@name(\"ФИО\", дательный, false, 1)")

	widthEntry := widget.NewEntry()
	if f.Width != 0 {
		widthEntry.SetText(fmt.Sprintf("%.2f", f.Width))
	}
	widthEntry.SetPlaceHolder("100.00")

	fontEntry := widget.NewEntry()
	fontEntry.SetText(f.Font)
	fontEntry.SetPlaceHolder("PT Sans")

	fontSizeEntry := widget.NewEntry()
	if f.FontSize != 0 {
		fontSizeEntry.SetText(fmt.Sprintf("%.2f", f.FontSize))
	}
	fontSizeEntry.SetPlaceHolder("14.00")

	// необязательные поля
	xEntry := widget.NewEntry()
	if f.X != 0 || f.Indent != "" {
		xEntry.SetText(fmt.Sprintf("%.2f", f.X))
	}
	xEntry.SetPlaceHolder("0.00")

	yEntry := widget.NewEntry()
	if f.Y != 0 || f.Indent != "" {
		yEntry.SetText(fmt.Sprintf("%.2f", f.Y))
	}
	yEntry.SetPlaceHolder("100.00")

	colorEntry := widget.NewEntry()
	colorEntry.SetText(f.Color)
	colorEntry.SetPlaceHolder("#000000")

	fontTypeEntry := widget.NewEntry()
	fontTypeEntry.SetText(f.FontType)
	fontTypeEntry.SetPlaceHolder("regular")

	alignEntry := widget.NewEntry()
	alignEntry.SetText(f.Align)
	alignEntry.SetPlaceHolder("center")

	placeAlignEntry := widget.NewEntry()
	placeAlignEntry.SetText(f.PlaceAlign)
	placeAlignEntry.SetPlaceHolder("top + left")

	leadingEntry := widget.NewEntry()
	if f.Leading != 0 {
		leadingEntry.SetText(fmt.Sprintf("%.2f", f.Leading))
	}
	leadingEntry.SetPlaceHolder("5.00")

	spacingEntry := widget.NewEntry()
	if f.Spacing != 0 {
		spacingEntry.SetText(fmt.Sprintf("%.2f", f.Spacing))
	}
	spacingEntry.SetPlaceHolder("100.00")

	trackingEntry := widget.NewEntry()
	if f.Tracking != 0 {
		trackingEntry.SetText(fmt.Sprintf("%.2f", f.Tracking))
	}
	trackingEntry.SetPlaceHolder("0.00")

	indentEntry := widget.NewEntry()
	indentEntry.SetText(f.Indent)
	indentEntry.SetPlaceHolder("")

	indentValEntry := widget.NewEntry()
	if f.IndentValue != 0 {
		indentValEntry.SetText(fmt.Sprintf("%.2f", f.IndentValue))
	}
	indentValEntry.SetPlaceHolder("0.00")

	fieldForm := widget.NewForm(
		// обязательные (5)
		widget.NewFormItem("Field ID *", idEntry),
		widget.NewFormItem("Value *", valEntry),
		widget.NewFormItem("Width (mm) *", widthEntry),
		widget.NewFormItem("Font *", fontEntry),
		widget.NewFormItem("Font Size (pt) *", fontSizeEntry),

		// необязательные (11)
		widget.NewFormItem("X (mm)", xEntry),
		widget.NewFormItem("Y (mm)", yEntry),
		widget.NewFormItem("Hex Color", colorEntry),
		widget.NewFormItem("Font Type", fontTypeEntry),
		widget.NewFormItem("Align", alignEntry),
		widget.NewFormItem("Place Align", placeAlignEntry),
		widget.NewFormItem("Leading (mm)", leadingEntry),
		widget.NewFormItem("Spacing (%)", spacingEntry),
		widget.NewFormItem("Tracking (mm)", trackingEntry),
		widget.NewFormItem("Indent ID", indentEntry),
		widget.NewFormItem("Indent Value (mm)", indentValEntry),
	)

	var item *widget.AccordionItem
	deleteBut := widget.NewButtonWithIcon("Delete Field", theme.DeleteIcon(), func() {
		c.fieldsAccordion.Remove(item)
		for i, fui := range c.fieldsUI {
			if fui.accordionItem == item {
				c.fieldsUI = append(c.fieldsUI[:i], c.fieldsUI[i+1:]...)
				break
			}
		}
		c.fieldsAccordion.Refresh()
	})

	fieldContainer := container.NewVBox(fieldForm, deleteBut)

	header := f.ID
	if header == "" {
		header = "new_field"
	}
	item = widget.NewAccordionItem(header, fieldContainer)

	idEntry.OnChanged = func(newID string) {
		if newID != "" {
			item.Title = newID
		} else {
			item.Title = "empty_id"
		}
		c.fieldsAccordion.Refresh()
	}

	c.fieldsUI = append(c.fieldsUI, &FieldUI{
		idEntry:         idEntry,
		valEntry:        valEntry,
		xEntry:          xEntry,
		yEntry:          yEntry,
		widthEntry:      widthEntry,
		fontEntry:       fontEntry,
		fontTypeEntry:   fontTypeEntry,
		fontSizeEntry:   fontSizeEntry,
		colorEntry:      colorEntry,
		alignEntry:      alignEntry,
		placeAlignEntry: placeAlignEntry,
		leadingEntry:    leadingEntry,
		spacingEntry:    spacingEntry,
		trackingEntry:   trackingEntry,
		indentEntry:     indentEntry,
		indentValEntry:  indentValEntry,
		accordionItem:   item,
	})

	c.fieldsAccordion.Append(item)
	c.fieldsAccordion.Refresh()
}
