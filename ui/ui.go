package ui

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/input"
)

type Element interface {
	Update(inputs *input.Input)

	SetPosition(x, y float32)
	SetDimensions(width, height float32)

	GetTransform() geometry.Transform

	GetChildren() []*child.Child2D
	GetTextBoxes() []*TextBox
}
