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

func NewUIButton(x, y, width, height float32, material *material.Material, config *configuration.EngineConfig) Button {
	button := Button{
		justClicked: false,
		colliding:   make(map[int]bool),
	}

	c := child.NewChild2D(config)
	c.AttachMesh(geometry.NewRectangle(width, height, config))
	c.AttachMaterial(material)
	c.AttachCollider(0, 0, width, height)
	c.X = x
	c.Y = y
	c.SetMouseFunc(button.MouseFunc)

	button.ElementChild = &c

	return button
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
