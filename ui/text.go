package ui

import (
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
)

type TextBox struct {
	Text    string
	textObj *v41.Text
	Font    string

	Scale float32

	X     float32
	Y     float32
	Color [3]float32
}

func (t *TextBox) Update() {
	t.textObj.SetString(t.Text)
	t.textObj.SetPosition(mgl32.Vec2{t.X, t.Y})
	t.textObj.SetScale(t.Scale)
	t.textObj.SetColor(mgl32.Vec3(t.Color))
	t.textObj.Draw()
}

func (t *TextBox) SetV41Text(textObj *v41.Text) {
	t.textObj = textObj
}
