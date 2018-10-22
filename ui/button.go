package ui

import (
	"rapidengine/child"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/input"
	"rapidengine/material"
)

type Button struct {
	ElementChild *child.Child2D

	width  float32
	height float32

	clickCallback func()
	justClicked   bool

	colliding map[int]bool
}

func NewUIButton(x, y, width, height float32, material *material.Material, config *configuration.EngineConfig) (Button, *child.Child2D) {
	button := Button{
		justClicked: false,
		colliding:   make(map[int]bool),
	}

	child := child.NewChild2D(config)
	child.AttachPrimitive(geometry.NewRectangle(width, height, config))
	child.AttachMaterial(material)
	child.AttachCollider(0, 0, width, height)
	//child.AttachShader(engine.ShaderControl.GetShader("color"))
	child.SetPosition(x, y)
	child.SetMouseFunc(button.MouseFunc)

	button.ElementChild = &child

	return button, &child
}

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

func (button *Button) SetClickCallback(f func()) {
	button.clickCallback = f
}

func (button *Button) MouseFunc(c bool) {
	button.colliding[0] = c
}
