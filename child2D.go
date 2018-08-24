package main

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Child2D struct {
	vertexArray *VertexArray
	numVertices int32

	shaderProgram uint32
	texture       uint32

	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4

	X float32
	Y float32
}

func NewChild2D() Child2D {
	return Child2D{
		modelMatrix:      mgl32.Ident4(),
		projectionMatrix: mgl32.Ortho2D(-1, 1, -1, 1),
	}
}

func (child2D *Child2D) PreRender(mainCamera Camera) {
	gl.BindVertexArray(child2D.vertexArray.id)
	gl.UseProgram(child2D.shaderProgram)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.shaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child2D.projectionMatrix[0],
	)

	gl.BindAttribLocation(child2D.shaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child2D.shaderProgram, 1, gl.Str("tex\x00"))
}

func (child2D *Child2D) Update() {
	child2D.modelMatrix = mgl32.Translate3D(-1, 1, 0)
	child2D.projectionMatrix = mgl32.Ortho2D(-1, 1, -1, 1)
}

func (child2D *Child2D) AttachTexture(path string, coords []float32) error {
	if child2D.vertexArray == nil {
		return errors.New("Cannot attach texture without VertexArray")
	}
	if child2D.shaderProgram == 0 {
		return errors.New("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child2D.vertexArray.id)
	gl.UseProgram(child2D.shaderProgram)

	child2D.vertexArray.AddVertexAttribute(coords, 1, 2)

	texture, err := NewTexture(path)
	if err != nil {
		return err
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	loc1 := gl.GetUniformLocation(child2D.shaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, int32(0))
	CheckError("IGNORE")

	child2D.texture = texture
	gl.BindVertexArray(0)
	return nil
}

func (child2D *Child2D) AttachVertexArray(vao *VertexArray, numVertices int32) {
	child2D.vertexArray = vao
	child2D.numVertices = numVertices
}

func (child2D *Child2D) AttachPrimitive(p Primitive) {
	child2D.AttachVertexArray(p.vao, p.numVertices)
}

func (child2D *Child2D) AttachShader(s uint32) {
	child2D.shaderProgram = s
}

func (child2D *Child2D) SetPosition(x, y float32) {
	child2D.X = x
	child2D.Y = y
}

func (child2D *Child2D) GetShaderProgram() uint32 {
	return child2D.shaderProgram
}

func (child2D *Child2D) GetVertexArray() *VertexArray {
	return child2D.vertexArray
}

func (child2D *Child2D) GetNumVertices() int32 {
	return child2D.numVertices
}
