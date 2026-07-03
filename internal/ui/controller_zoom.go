package ui

import "fyne.io/fyne/v2"

func (c *UIController) ZoomIn() {
	if c.zoomLevel < 4.0 {
		c.zoomLevel *= 1.2
		c.updateZoom()
	}
}

func (c *UIController) ZoomOut() {
	viewportSize := c.imageScroll.Size()
	minWidth := float32(400)
	minHeight := float32(300)
	if viewportSize.Width > 0 {
		minWidth = viewportSize.Width
	}
	if viewportSize.Height > 0 {
		minHeight = viewportSize.Height
	}

	nextZoomLevel := c.zoomLevel / 1.2
	nextWidth := 1000 * nextZoomLevel
	nextHeight := 1000 * nextZoomLevel

	// разрешаем отдаление, если хотя бы по одной оси картинка все еще больше вьюпорта
	if nextWidth >= minWidth || nextHeight >= minHeight {
		c.zoomLevel = nextZoomLevel
		c.updateZoom()
	}
}

func (c *UIController) updateZoom() {
	// ограничение приближения, максимум 4.0
	if c.zoomLevel >= 4.0 {
		c.zoomInBut.Disable()
	} else {
		c.zoomInBut.Enable()
	}

	// проверяем, можно ли будет скроллить при следующем шаге отдаления
	viewportSize := c.imageScroll.Size()
	minWidth := float32(400)
	minHeight := float32(300)
	if viewportSize.Width > 0 {
		minWidth = viewportSize.Width
	}
	if viewportSize.Height > 0 {
		minHeight = viewportSize.Height
	}

	nextZoomLevel := c.zoomLevel / 1.2
	nextWidth := 1000 * nextZoomLevel
	nextHeight := 1000 * nextZoomLevel

	// если при следующем шаге картинка полностью влезет по обеим осям, выключаем кнопку отдаления
	if nextWidth < minWidth && nextHeight < minHeight {
		c.zoomOutBut.Disable()
	} else {
		c.zoomOutBut.Enable()
	}

	w := 1000 * c.zoomLevel
	h := 1000 * c.zoomLevel
	c.imageLayout.size = fyne.NewSize(w, h)
	c.imageScroll.Content.Refresh()
	c.imageScroll.Refresh()
}
