package cmd

import (
	"rapidengine/input"
	"rapidengine/material"
	"rapidengine/ui"
)

type UIControl struct {
	Buttons map[string]*ui.Button

	engine *Engine
}

func NewUIControl() UIControl {
	return UIControl{
		Buttons: make(map[string]*ui.Button),
	}
}

func (uiControl *UIControl) Initialize(engine *Engine) {
	uiControl.engine = engine
}

func (uiControl *UIControl) Update(inputs *input.Input) {
	for _, button := range uiControl.Buttons {
		button.Update(inputs)
	}
}

func (uiControl *UIControl) AddButton(button *ui.Button, id string) {
	uiControl.engine.Instance(button.ElementChild)
	uiControl.engine.CollisionControl.CreateMouseCollision(button.ElementChild)
	uiControl.engine.TextControl.AddTextBox(button.TextBx)

	uiControl.Buttons[id] = button
}

func (uiControl *UIControl) NewUIButton(
	x, y,
	width, height float32,
	material *material.Material,
) ui.Button {
	button := ui.NewUIButton(x, y, width, height, material, uiControl.engine.Config)

	button.ElementChild.AttachShader(uiControl.engine.ShaderControl.GetShader("color"))

	return button
}
