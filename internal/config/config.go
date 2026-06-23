// Package config ответчает за загрузку и хранение конфигураций шаблонов
package config

type Field struct {
	ID          string  `json:"id"`           // ID блока
	Value       string  `json:"value"`        // Сама строка
	X           float64 `json:"x"`            // Координаты x в mm
	Y           float64 `json:"y"`            // Координаты y в mm
	Width       float64 `json:"width"`        // Ширина в mm
	Indent      string  `json:"indent"`       // Есть ли отступ от блока выше динамический. Указывается ID блока относительно которого отсутп идет. Стандартное значение "". Если стоит ID, то Y не учитывается
	IndentValue float64 `json:"indent_value"` // Если есть Indent, то тут указывается отступ в em. Если Indent пустое, то это поле игнорируется
	Align       string  `json:"align"`        // Вырванивание текста внутри блока: center, left, right + top/bottom
	PlaceAlign  string  `json:"place_align"`  // Выравнивание самого блока текст: center, left, right + top/bottom
	Leading     float64 `json:"leading"`      // Межстрочный интервал (в em)
	Spacing     float64 `json:"spacing"`      // Пробелы между словами (в %, по стандарту 100)
	Tracking    float64 `json:"tracking"`     // Расстояние между буквами в мм
	Color       string  `json:"color"`        // HEX кодировка цвета шрифта
	Font        string  `json:"font"`         // Название шрифта
	FontType    string  `json:"font_type"`    // Тип шрифта: regular, thin, bold, black e.t.c.
	FontSize    float64 `json:"font_size"`    // Размер шрифта в pt
}

type Background struct {
	BackgroundImage string  `json:"background_image"` // Подложка (pdf, png, svg e.t.c.)
	BgScale         float64 `json:"bg_scale"`         // Масштаб фона относительно центра в % (100% по размеру страницы, т.е. макс размер, чтобы не вышло за рамки). Необходимо, чтобы слегка увеличить подложку для печати без белых рамок
	BgDx            float64 `json:"bg_dx"`            // Смещение подложки от центра по ширине в mm
	BgDy            float64 `json:"bg_dy"`            // Смещение подложки от центра по высоте в mm
}

type TemplateConfig struct {
	Title  string     `json:"title"`      // Название шаблона
	Width  float64    `json:"width"`      // Ширина страницы в mm
	Height float64    `json:"height"`     // Высота страницы в mm
	Fields []Field    `json:"fields"`     // Срез полей
	Bg     Background `json:"background"` // Подложка
}
