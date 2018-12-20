package ui

import (
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/input"
)

type RootElement struct {
	// State: none

	// Tree
	Elements []Element

	// Position
	transform geometry.Transform
}

func NewRootElement(width, height float32) RootElement {
	return RootElement{
		transform: geometry.NewTransform(0, 0, 0, width, height, 0),
	}
}

func (re *RootElement) Update(inputs *input.Input) {

}

func (re *RootElement) GetChildren() []child.Child {
	children := []child.Child{}
	for _, e := range re.Elements {
		children = append(children, e.GetChildren()...)
	}
	return children
}

func (re *RootElement) GetElements() []Element {
	return re.Elements
}

func (re *RootElement) GetTransform() *geometry.Transform {
	return &re.transform
}

func (re *RootElement) GetTextBoxes() []*TextBox {
	return []*TextBox{}
}
