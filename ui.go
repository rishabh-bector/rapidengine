package rapidengine

import (
	"rapidengine/input"
)

type UIControl struct {
	Buttons map[string]*UIButton
}

func NewUIControl() UIControl {
	return UIControl{
		Buttons: make(map[string]*UIButton),
	}
}

func (ui *UIControl) Update(inputs *input.Input) {
	for _, button := range ui.Buttons {
		button.Update(inputs)
	}
}

func (ui *UIControl) AddButton(button *UIButton, id string) {
	ui.Buttons[id] = button
}

type UIButton struct {
	ElementChild *Child2D

	width  float32
	height float32

	clickCallback func()
	justClicked   bool

	colliding map[int]bool
}

func NewUIButton(x, y, width, height float32, material *Material, engine *Engine) UIButton {
	button := UIButton{
		justClicked: false,
		colliding:   make(map[int]bool),
	}

	child := NewChild2D(&engine.Config, &engine.CollisionControl)
	child.AttachPrimitive(NewRectangle(width, height, &engine.Config))
	child.AttachMaterial(material)
	child.AttachCollider(0, 0, width, height)
	child.AttachShader(engine.ShaderControl.GetShader("color"))
	child.SetPosition(x, y)
	child.SetMouseFunc(button.MouseFunc)
	engine.Instance(&child)

	engine.CollisionControl.CreateMouseCollision(&child)
	button.ElementChild = &child

	return button
}

func (button *UIButton) Update(inputs *input.Input) {
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

func (button *UIButton) SetClickCallback(f func()) {
	button.clickCallback = f
}

func (button *UIButton) MouseFunc(c bool) {
	button.colliding[0] = c
}
