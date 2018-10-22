package control

import (
	"fmt"
	"os"
	"rapidengine/configuration"
	"rapidengine/ui"

	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/math/fixed"
)

type TextControl struct {
	Texts map[int]*ui.TextBox
	Fonts map[string]*v41.Font

	numTexts int

	config *configuration.EngineConfig
}

func NewTextControl(config *configuration.EngineConfig) TextControl {
	return TextControl{
		Texts:    make(map[int]*ui.TextBox),
		Fonts:    make(map[string]*v41.Font),
		numTexts: 0,
		config:   config,
	}
}

func (tc *TextControl) Update() {
	for _, textbox := range tc.Texts {
		textbox.Update()
	}
}

func (tc *TextControl) NewTextBox(text string, font string, x, y, scale float32, color [3]float32) *ui.TextBox {
	t := v41.NewText(tc.Fonts[font], 0.2, 10)
	t.SetString(text)
	t.SetColor(mgl32.Vec3{1, 1, 1})
	t.AddScale(scale)

	textbox := &ui.TextBox{
		Text:  text,
		Font:  font,
		Color: color,
		X:     x,
		Y:     y,
		Scale: scale,
	}
	textbox.SetV41Text(t)

	return textbox
}

func (tc *TextControl) AddTextBox(tb *ui.TextBox) {
	tc.Texts[tc.numTexts] = tb
	tc.numTexts++
}

func (tc *TextControl) LoadFont(path string, name string, scale float32, offset int) {
	var font *v41.Font
	config, err := gltext.LoadTruetypeFontConfig("fontconfigs", name)
	if err == nil {
		font, err = v41.NewFont(config)
		if err != nil {
			panic(err)
		}
		fmt.Println("Font loaded from disk...")
	} else {
		fd, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		// Japanese character ranges
		// http://www.rikai.com/library/kanjitables/kanji_codes.unicode.shtml
		runeRanges := make(gltext.RuneRanges, 0)
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 32, High: 128})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3000, High: 0x3030})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3040, High: 0x309f})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x30a0, High: 0x30ff})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x4e00, High: 0x9faf})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0xff00, High: 0xffef})

		scale := fixed.Int26_6(scale)
		runesPerRow := fixed.Int26_6(128)
		config, err = gltext.NewTruetypeFontConfig(fd, scale, runeRanges, runesPerRow, fixed.Int26_6(offset))
		if err != nil {
			panic(err)
		}
		err = config.Save("fontconfigs", name)
		if err != nil {
			panic(err)
		}
		font, err = v41.NewFont(config)
		if err != nil {
			panic(err)
		}
	}

	font.ResizeWindow(float32(tc.config.ScreenWidth), float32(tc.config.ScreenHeight))

	tc.Fonts[name] = font

}
