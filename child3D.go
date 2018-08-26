package rapidengine

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Child3D struct {
	vertexArray *VertexArray
	numVertices int32

	shaderProgram uint32
	texture       uint32

	modelMatrix      mgl32.Mat4
	projectionMatrix mgl32.Mat4

	Config *EngineConfig
}

func NewChild3D() Child3D {
	return Child3D{
		modelMatrix:      mgl32.Ident4(),
		projectionMatrix: mgl32.Ident4(),
	}
}

func (child3D *Child3D) PreRender(mainCamera Camera3D) {
	gl.BindVertexArray(child3D.vertexArray.id)
	gl.UseProgram(child3D.shaderProgram)

	child3D.modelMatrix = mgl32.Ident4()

	child3D.projectionMatrix = mgl32.Perspective(
		mgl32.DegToRad(45),
		float32(child3D.Config.ScreenWidth)/float32(child3D.Config.ScreenHeight),
		0.1, 100,
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child3D.modelMatrix[0],
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("viewMtx\x00")),
		1, false, &mainCamera.View[0],
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child3D.shaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child3D.projectionMatrix[0],
	)

	gl.BindAttribLocation(child3D.shaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child3D.shaderProgram, 1, gl.Str("tex\x00"))
}

func (child3D *Child3D) AttachTexture(path string, coords []float32) error {
	if child3D.vertexArray == nil {
		return errors.New("Cannot attach texture without VertexArray")
	}
	if child3D.shaderProgram == 0 {
		return errors.New("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child3D.vertexArray.id)
	gl.UseProgram(child3D.shaderProgram)

	child3D.vertexArray.AddVertexAttribute(coords, 1, 2)

	texture, err := NewTexture(path)
	if err != nil {
		return err
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	loc1 := gl.GetUniformLocation(child3D.shaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, int32(0))
	CheckError("IGNORE")

	child3D.texture = texture
	gl.BindVertexArray(0)
	return nil
}

func (child3D *Child3D) AttachVertexArray(vao *VertexArray, numVertices int32) {
	child3D.vertexArray = vao
	child3D.numVertices = numVertices
}

func (child3D *Child3D) AttachPrimitive(p Primitive) {
	child3D.AttachVertexArray(p.vao, p.numVertices)
}

func (child3D *Child3D) AttachShader(s uint32) {
	child3D.shaderProgram = s
}
