package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Renderer struct {
	window        *glfw.Window
	shaderProgram uint32
	children      []Child
}

func (renderer *Renderer) RenderChildren() {
	for _, child := range renderer.children {
		child.Render()
	}
}

func (renderer *Renderer) Register(child Child) {
	child.PreRender()
	renderer.children = append(renderer.children, child)
}
