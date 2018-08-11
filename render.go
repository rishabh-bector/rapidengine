package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Renderer struct {
	window        *glfw.Window
	shaderProgram uint32
}

func (renderer *Renderer) DrawElements(vertexArray VertexArray, numElements int32) {
	gl.LinkProgram(renderer.shaderProgram)
	gl.UseProgram(renderer.shaderProgram)
	gl.BindVertexArray(vertexArray.id)
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, numElements, gl.UNSIGNED_INT, gl.PtrOffset(0))
}
