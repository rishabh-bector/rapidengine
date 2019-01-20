package ui

import (
	"rapidengine/child"
	"rapidengine/input"
)

type Element interface {
	Update(inputs *input.Input)

	// Get all children (for scene instancing)
	GetChildren() []child.Child

	// UI Tree
	InstanceElement(Element)
	GetElements() []Element

	// Components
	GetTransform() *UITransform
	GetTextBox() *TextBox
}

type UITransform struct {
	X float32
	Y float32

	PadX float32
	PadY float32

	SX float32
	SY float32

	Parent Element
}

func (t *UITransform) UpdateChild(c *child.Child2D) {
	c.X = t.X
	c.Y = t.Y
	c.ScaleX = t.SX
	c.ScaleY = t.SY
}

func (t *UITransform) AlignToParent() {
	t.X = t.Parent.GetTransform().X + t.PadX
	t.Y = t.Parent.GetTransform().Y + t.PadY
}
