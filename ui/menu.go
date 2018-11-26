package ui

import (
	"rapidengine/child"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/input"
)

type Menu struct {
	elements []Element

	BackChild *child.Child2D

	transform geometry.Transform
}

func NewMenu(config *configuration.EngineConfig) Menu {
	return Menu{
		transform: geometry.NewTransform(0, 0, 0, 200, 100, 0),
	}
}

func (m *Menu) Initialize() {
	m.BackChild.SetPosition(m.transform.X, m.transform.Y)

	m.BackChild.ScaleX = m.transform.SX
	m.BackChild.ScaleY = m.transform.SY
}

func (m *Menu) AddElement(e Element) {
	m.elements = append(m.elements, e)
}

//  --------------------------------------------------
//  Interface
//  --------------------------------------------------

func (m *Menu) Update(inputs *input.Input) {

}

func (m *Menu) SetPosition(x, y float32) {
	m.transform.X = x
	m.transform.Y = y
	m.Initialize()
}

func (m *Menu) SetDimensions(width, height float32) {
	m.transform.SX = width
	m.transform.SY = height
	m.Initialize()
}

func (m *Menu) GetTransform() geometry.Transform {
	return m.transform
}

func (m *Menu) GetChildren() []*child.Child2D {
	children := []*child.Child2D{m.BackChild}
	for _, e := range m.elements {
		children = append(children, e.GetChildren()...)
	}
	return children
}

func (m *Menu) GetTextBoxes() []*TextBox {
	tbs := []*TextBox{}
	for _, e := range m.elements {
		tbs = append(tbs, e.GetTextBoxes()...)
	}
	return tbs
}
