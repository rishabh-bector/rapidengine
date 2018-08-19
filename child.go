package main

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Child struct {
	vertexArray   *VertexArray
	numVertices   int32
	shaderProgram uint32
	texture       uint32

	modelMatrix      mgl32.Mat4
	viewMatrix       mgl32.Mat4
	projectionMatrix mgl32.Mat4
}

func NewChild() Child {
	return Child{
		modelMatrix:      mgl32.Ident4(),
		viewMatrix:       mgl32.Ident4(),
		projectionMatrix: mgl32.Ident4(),
	}
}

func (child *Child) PreRender() {
	gl.BindVertexArray(child.vertexArray.id)
	gl.UseProgram(child.shaderProgram)

	child.modelMatrix = mgl32.HomogRotate3D(float32(glfw.GetTime())*mgl32.DegToRad(-55), mgl32.Vec3{0, 1, 0})

	child.viewMatrix = mgl32.Translate3D(-0.5, 0, -3)
	child.projectionMatrix = mgl32.Perspective(
		mgl32.DegToRad(45),
		float32(ScreenWidth)/float32(ScreenHeight),
		0.1, 100,
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child.modelMatrix[0],
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, &child.viewMatrix[0],
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child.shaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child.projectionMatrix[0],
	)

	gl.BindAttribLocation(child.shaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child.shaderProgram, 1, gl.Str("tex\x00"))
}

func (child *Child) AttachTexture(path string, coords []float32) error {
	if child.vertexArray == nil {
		return errors.New("Cannot attach texture without VertexArray")
	}
	if child.shaderProgram == 0 {
		return errors.New("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child.vertexArray.id)
	gl.UseProgram(child.shaderProgram)

	child.vertexArray.AddVertexAttribute(coords, 1, 2)

	texture, err := NewTexture(path)
	if err != nil {
		return err
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	loc1 := gl.GetUniformLocation(child.shaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, int32(0))
	CheckError("IGNORE")

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
