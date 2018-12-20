package ui

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/input"
)

type Button struct {
	// State
	ButtonChild   *child.Child2D
	clickCallback func()
	justClicked   bool
	colliding     map[int]bool
	TextBx        *TextBox

	// Tree
	Parent   Element
	Elements []Element

	// Position
	transform geometry.Transform
	AlignX    int
	AlignY    int
}

func NewUIButton(x, y, width, height float32) Button {
	button := Button{
		justClicked: false,
		colliding:   make(map[int]bool),
		TextBx:      nil,

		transform: geometry.NewTransform(x, y, 0, width, height, 0),
		AlignX:    -1,
		AlignY:    -1,
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
	button.ButtonChild.X = button.transform.X
	button.ButtonChild.Y = button.transform.Y

	button.ButtonChild.ScaleX = button.transform.SX
	button.ButtonChild.ScaleY = button.transform.SY

	if button.TextBx != nil {
		button.TextBx.X = button.ButtonChild.X + button.GetTransform().SX/2
		button.TextBx.Y = button.ButtonChild.Y + button.GetTransform().SY/2
	}

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

func (button *Button) GetTransform() *geometry.Transform {
	return &button.transform
}

func (button *Button) GetChildren() []child.Child {
	return []child.Child{button.ButtonChild}
}

func (button *Button) GetTextBoxes() []*TextBox {
	return []*TextBox{button.TextBx}
}

func (button *Button) GetElements() []Element {
	return button.Elements
}
