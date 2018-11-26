package ui

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/input"
)

type Button struct {
	ButtonChild *child.Child2D

	TextBx *TextBox

	transform geometry.Transform

	clickCallback func()
	justClicked   bool

	colliding map[int]bool
}

func NewUIButton(x, y, width, height float32) Button {
	button := Button{
		justClicked: false,
		colliding:   make(map[int]bool),
		TextBx:      nil,
		transform:   geometry.NewTransform(x, y, 0, width, height, 0),
	}

	return button
}

func (button *Button) Initialize() {
	button.ButtonChild.AttachCollider(
		0, 0,
		button.transform.SX,
		button.transform.SY,
	)
	button.ButtonChild.SetMouseFunc(button.MouseFunc)
	button.ButtonChild.Static = true
	button.SetPosition(button.transform.X, button.transform.Y)
	button.SetDimensions(button.transform.SX, button.transform.SY)
}

func (button *Button) SetClickCallback(f func()) {
	button.clickCallback = f
}

func (button *Button) AttachText(tb *TextBox) {
	button.TextBx = tb
	button.Initialize()
}

func (button *Button) MouseFunc(c bool) {
	button.colliding[0] = c
}

func (button *Button) Block() {
	button.justClicked = true
}

//  --------------------------------------------------
//  Interface
//  --------------------------------------------------

func (button *Button) Update(inputs *input.Input) {
	if button.colliding[0] {
		if inputs.LeftMouseButton {
			if !button.justClicked {
				button.clickCallback()
				button.justClicked = true
			}
		} else {
			button.justClicked = false
		}
	}
}

func (button *Button) SetPosition(x, y float32) {
	button.ButtonChild.X = x
	button.ButtonChild.Y = y

	if button.TextBx != nil {
		button.TextBx.X = button.ButtonChild.X + button.GetTransform().SX/2
		button.TextBx.Y = button.ButtonChild.Y + button.GetTransform().SY/2
	}
}

func (button *Button) SetDimensions(width, height float32) {
	button.ButtonChild.ScaleX = width
	button.ButtonChild.ScaleY = height
}

func (button *Button) GetTransform() geometry.Transform {
	return button.transform
}

func (button *Button) GetChildren() []*child.Child2D {
	return []*child.Child2D{button.ButtonChild}
}

func (button *Button) GetTextBoxes() []*TextBox {
	return []*TextBox{button.TextBx}
}
