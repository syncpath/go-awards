package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ZoomLayout struct {
	size fyne.Size
}

func (z *ZoomLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return z.size
}

func (z *ZoomLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, obj := range objects {
		obj.Resize(size)
		obj.Move(fyne.NewPos(0, 0))
	}
}

type ReadOnlyEntry struct {
	widget.Entry
}

func NewReadOnlyEntry() *ReadOnlyEntry {
	e := &ReadOnlyEntry{}
	e.MultiLine = true
	e.ExtendBaseWidget(e)
	return e
}

func (e *ReadOnlyEntry) TypedRune(r rune) {}
func (e *ReadOnlyEntry) TypedKey(key *fyne.KeyEvent) {
	// разрешаем только клавиши навигации
	switch key.Name {
	case fyne.KeyUp, fyne.KeyDown, fyne.KeyLeft, fyne.KeyRight,
		fyne.KeyHome, fyne.KeyEnd, fyne.KeyPageUp, fyne.KeyPageDown:
		e.Entry.TypedKey(key)
	}
}

func (e *ReadOnlyEntry) TypedShortcut(shortcut fyne.Shortcut) {
	// разрешаем только копирование и выделение всего текста
	switch shortcut.(type) {
	case *fyne.ShortcutCopy, *fyne.ShortcutSelectAll:
		e.Entry.TypedShortcut(shortcut)
	}
}
