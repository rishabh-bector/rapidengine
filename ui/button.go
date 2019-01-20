package ui

import (
	"rapidengine/child"
	"rapidengine/input"
)

type Button struct {
	// State
	ButtonChild   *child.Child2D
	clickCallback func()
	justClicked   bool
	colliding     map[int]bool

	TextBx  *TextBox
	TextPad float32

	// Tree
	Elements []Element

	// Position
	Transform *UITransform
}

func NewUIButton(x, y, width, height float32) Button {
	button := Button{
		justClicked: false,
		colliding:   make(map[int]bool),
		TextBx:      nil,

		Transform: &UITransform{
			X:    x,
			Y:    y,
			PadX: x,
			PadY: y,
			SX:   width,
			SY:   height,
		},
	}

	return button
}

func (button *Button) Initialize() {
	button.ButtonChild.AttachCollider(
		0, 0,
		button.Transform.SX,
		button.Transform.SY,
	)
	button.ButtonChild.SetMouseFunc(button.MouseFunc)
	button.ButtonChild.Static = true
}

func (button *Button) SetClickCallback(f func()) {
	button.clickCallback = f
}

func (button *Button) AttachText(tb *TextBox) {
	button.TextBx = tb
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
	// Align in parent
	button.Transform.AlignToParent()

	// Update child transform
	button.Transform.UpdateChild(button.ButtonChild)

	// Update textbox transform
	button.TextBx.X = button.ButtonChild.X + button.GetTransform().SX/2
	button.TextBx.Y = button.ButtonChild.Y + button.GetTransform().SY/2
	if button.TextBx.Text != "" {
		//button.Transform.SX = float32(button.TextBx.GetLength())*20 + button.TextPad
		//button.ButtonChild.Collider.Width = button.Transform.SX
	}

	// Collision logic
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

	// Tree
	for _, e := range button.Elements {
		e.Update(inputs)
	}
}

func (button *Button) InstanceElement(e Element) {
	button.Elements = append(button.Elements, e)
	e.GetTransform().Parent = button
}

func (button *Button) GetElements() []Element {
	elements := []Element{}
	for _, e := range button.Elements {
		elements = append(elements, e)
		elements = append(elements, e.GetElements()...)
	}
	return elements
}

func (button *Button) GetTransform() *UITransform {
	return button.Transform
}

func (button *Button) GetChildren() []child.Child {
	return []child.Child{button.ButtonChild}
}

func (button *Button) GetTextBox() *TextBox {
	return button.TextBx
}
