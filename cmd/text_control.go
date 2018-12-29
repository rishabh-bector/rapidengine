package cmd

import (
	"fmt"
	"os"
	"rapidengine/configuration"
	"rapidengine/state"
	"rapidengine/ui"

	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/math/fixed"
)

type TextControl struct {
	Fonts map[string]*v41.Font

	engine *Engine
}

func NewTextControl(config *configuration.EngineConfig) TextControl {
	return TextControl{
		Fonts: make(map[string]*v41.Font),
	}
}

func (tc *TextControl) Initialize(engine *Engine) {
	tc.engine = engine
}

func (tc *TextControl) Update() {
	for _, t := range tc.engine.SceneControl.GetCurrentTexts() {
		t.Update(tc.engine.Config)
	}
	state.BoundTexture0 = 999
}

func (tc *TextControl) NewTextBox(text string, font string, x, y, scale float32, color [3]float32) *ui.TextBox {
	t := v41.NewText(tc.Fonts[font], 0.2, 10)
	t.SetString(text)
	t.SetColor(mgl32.Vec3{1, 1, 1})
	t.AddScale(scale)

	textbox := &ui.TextBox{
		Text:  text,
		Font:  font,
		Color: [3]float32{color[0] / 255, color[1] / 255, color[2] / 255},
		X:     x,
		Y:     y,
		Scale: scale,
	}
	textbox.SetV41Text(t)

	return textbox
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

	font.ResizeWindow(float32(tc.engine.Config.ScreenWidth), float32(tc.engine.Config.ScreenHeight))

	tc.Fonts[name] = font

}
