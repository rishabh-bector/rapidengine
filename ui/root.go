package ui

import (
	"rapidengine/child"
	"rapidengine/input"
)

type RootElement struct {
	// State: none

	// Tree
	Elements []Element

	// Position
	Transform UITransform
}

func NewRootElement(width, height float32) RootElement {
	return RootElement{
		Transform: UITransform{
			SX: width,
			SY: width,
		},
	}
}

func (re *RootElement) Update(inputs *input.Input) {
	for _, e := range re.Elements {
		e.Update(inputs)
	}
}

// GetElements returns all the elements in the tree
func (re *RootElement) GetElements() []Element {
	elements := []Element{}

	for _, e := range re.Elements {
		elements = append(elements, e.GetElements()...)
		elements = append(elements, e)
	}

	return elements
}

func (re *RootElement) InstanceElement(element Element) {
	re.Elements = append(re.Elements, element)
}

func (re *RootElement) GetChildren() []child.Child {
	children := []child.Child{}
	for _, e := range re.Elements {
		children = append(children, e.GetChildren()...)
	}
	return children
}

func (re *RootElement) GetTransform() *UITransform {
	return &re.Transform
}

func (re *RootElement) GetTextBoxes() []*TextBox {
	return []*TextBox{}
}
