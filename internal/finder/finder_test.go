package finder_test

import (
	"reflect"
	"testing"

	"github.com/syncpath/go-awards/internal/finder"
)

func TestFindFields(t *testing.T) {
	tests := []struct {
		name string
		data finder.PageData
		want finder.PageInfo
	}{
		{
			"Исправные данные",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						X:        50.0,
						Y:        30.0,
						S:        "Готовый текст",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMin:     0.0,
				YMax:     100.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        17.64,
						Y:        21.24,
						Text:     "Готовый текст",
						Leading:  0,
					},
				},
			},
		},
		{
			"Некоторые данные пропущены",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						S:        "Готовый текст",
					},
				},
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        0,
						Y:        21.24,
						Text:     "Готовый текст",
						Leading:  0,
					},
				},
			},
		},
		{
			"Два текста, маленький внутри большого",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        21.0,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        21.24,
						Text:     "Готовый текст!м",
						Leading:  0,
					},
				},
			},
		},
		{
			"Два текста, разный Y",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        31.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        21.0,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        20.88,
						Text:     "Готовый текст!",
						Leading:  0,
					},
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.41,
						Y:        21.24,
						Text:     "м",
						Leading:  -3.1,
					},
				},
			},
		},
		{
			"Два текста, разница в иксах меньше 2em (склейка)",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        47.0,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        21.24,
						Text:     "Готовый текст!м",
						Leading:  0,
					},
				},
			},
		},
		{
			"Два текста, разный Y",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        31.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        21.0,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        20.88,
						Text:     "Готовый текст!",
						Leading:  0,
					},
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.41,
						Y:        21.24,
						Text:     "м",
						Leading:  -3.1,
					},
				},
			},
		},
		{
			"Два текста, разница в иксах больше 2em (без склейки)",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        76.70,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        21.24,
						Text:     "Готовый текст!",
						Leading:  0,
					},
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        27.06,
						Y:        21.24,
						Text:     "м",
						Leading:  0,
					},
				},
			},
		},
		{
			"Два текста, разный Y",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        31.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        21.0,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        20.88,
						Text:     "Готовый текст!",
						Leading:  0,
					},
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.41,
						Y:        21.24,
						Text:     "м",
						Leading:  -3.1,
					},
				},
			},
		},
		{
			"Два текста, разница шрифтах",
			finder.PageData{
				Texts: []finder.Text{
					{
						Font:     "funnyRegular",
						FontSize: 14.0,
						Y:        30.0,
						X:        20.0,
						S:        "Готовый текст!",
					},
					{
						Font:     "funnyBold",
						FontSize: 14.0,
						Y:        30.0,
						X:        76.70,
						S:        "м",
					},
				},
				XMin:     0.0,
				XMax:     100.0,
				YMax:     100.0,
				YMin:     0.0,
				CoefsMap: map[string]float64{"funnyRegular": 0.7, "funnyBold": 0.7},
			},
			finder.PageInfo{
				Height: 35.28,
				Width:  35.28,
				Texts: []finder.TextField{
					{
						Font:     "funnyRegular",
						FontSize: 14,
						X:        7.06,
						Y:        21.24,
						Text:     "Готовый текст!",
						Leading:  0,
					},
					{
						Font:     "funnyBold",
						FontSize: 14,
						X:        27.06,
						Y:        21.24,
						Text:     "м",
						Leading:  0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := finder.FindFields(test.data)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("получил: %v, ожидал: %v", got, test.want)
			}
		})
	}
}

func TestFindFontsCoefs(t *testing.T) {
	tests := []struct {
		name      string
		capHeight float64
		want      float64
	}{
		{"ветка /1000 внутри окна", 700, 0.7},
		{"ветка /1000, ровно нижняя граница 600", 600, 0.6},
		{"ветка /1000, ровно верхняя граница 800", 800, 0.8},
		{"ветка /1000 ниже окна", 500, 0.7},
		{"ветка /1000 выше окна", 900, 0.7},
		{"ровно 1000 строгое > не срабатывает", 1000, 0.7},
		{"ветка /2048 внутри окна", 1400, 0.68},
		{"ветка /2048 ниже окна", 1100, 0.7},
		{"ветка /2048 выше окна", 1700, 0.7},
		{"capHeight = 0", 0, 0.7},
		{"capHeight отрицательный", -5, 0.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := finder.FindFontsCoefs(map[string]float64{"x": tt.capHeight})["x"]
			if got != tt.want {
				t.Errorf("получил: %v, ожидал: %v", got, tt.want)
			}
		})
	}
}

func TestCleanFontName(t *testing.T) {
	tests := []struct {
		name string
		font string
		want string
	}{
		{"префикс + суффикс", "ABCDEF+Funny-Bold", "Funny-Bold"},
		{"только суффикс через дефис", "Helvetica-Bold", "Helvetica-Bold"},
		{"верхний регистр", "FUNNYREGULAR", "FUNNYREGULAR"},
		{"несколько дефисов", "Foo-Bar-Baz", "Foo-Bar-Baz"},
		{"чистый", "arial", "arial"},
		{"без плюса и дефиса", "FunnyRegular", "FunnyRegular"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := finder.CleanFontName(tt.font)
			if got != tt.want {
				t.Errorf("получил: %q, ожидал: %q", got, tt.want)
			}
		})
	}
}
