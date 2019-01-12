package ui

import (
	"rapidengine/configuration"

	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
)

type TextBox struct {
	Text    string
	textObj *v41.Text
	Font    string

	Scale float32

	X float32
	Y float32

	Color [3]float32
}

func (t *TextBox) Update(config *configuration.EngineConfig) {
	t.textObj.SetString(t.Text)
	t.textObj.SetPosition(mgl32.Vec2{
		t.X - float32(config.ScreenWidth/2),
		t.Y - float32(config.ScreenHeight/2),
	})
	t.textObj.SetScale(t.Scale)
	t.textObj.SetColor(mgl32.Vec3(t.Color))
	t.textObj.Draw()
}

func (t *TextBox) SetV41Text(textObj *v41.Text) {
	t.textObj = textObj
}

func (t *TextBox) GetLength() int {
	return t.textObj.GetLength()
}
