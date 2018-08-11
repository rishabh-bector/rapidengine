package main

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Child struct {
	vertexArray   *VertexArray
	numVertices   int32
	shaderProgram uint32
	texture       uint32
}

func NewChild() Child {
	return Child{nil, 0, 0, 0}
}

func (child *Child) PreRender() {
	gl.BindAttribLocation(child.shaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child.shaderProgram, 1, gl.Str("tex\x00"))
}

func (child *Child) Render() {
	gl.LinkProgram(child.shaderProgram)
	gl.UseProgram(child.shaderProgram)
	gl.BindVertexArray(child.vertexArray.id)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.DrawElements(gl.TRIANGLES, child.numVertices, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (child *Child) AttachTexture(path string, coords []float32) error {
	if child.vertexArray == nil {
		return errors.New("Cannot attach texture without VertexArray")
	}
	gl.BindVertexArray(child.vertexArray.id)
	child.vertexArray.AddVertexAttribute(coords, 1, 2)
	texture, err := NewTexture(path)
	if err != nil {
		return err
	}
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	loc1 := gl.GetUniformLocation(child.shaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, 0)
	child.texture = texture
	gl.BindVertexArray(0)
	return nil
}

func (child *Child) AttachVertexArray(vao *VertexArray, numVertices int32) {
	child.vertexArray = vao
	child.numVertices = numVertices
}

func (child *Child) AttachPrimitive(p Primitive) {
	child.AttachVertexArray(p.vao, p.numVertices)
}

func (child *Child) AttachShader(s uint32) {
	child.shaderProgram = s
}
