package ui

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/input"
)

type Element interface {
	Update(inputs *input.Input)

	// Get all children (for scene instancing)
	GetChildren() []child.Child

	GetTransform() *geometry.Transform
	GetTextBoxes() []*TextBox

	GetElements() []Element
}
