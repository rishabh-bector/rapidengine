package rapidengine

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Child2D struct {
	VertexArray *VertexArray
	NumVertices int32

	Primitive string

	ShaderProgram uint32
	Texture       uint32

	ModelMatrix      mgl32.Mat4
	ProjectionMatrix mgl32.Mat4

	X float32
	Y float32

	Config *EngineConfig
}

func NewChild2D(config *EngineConfig) Child2D {
	return Child2D{
		ModelMatrix:      mgl32.Ident4(),
		ProjectionMatrix: mgl32.Ortho2D(-1, 1, -1, 1),
		Config:           config,
	}
}

func (child2D *Child2D) PreRender(mainCamera Camera) {
	gl.BindVertexArray(child2D.VertexArray.id)
	gl.UseProgram(child2D.ShaderProgram)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.ShaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.ModelMatrix[0],
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.ShaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.ShaderProgram, gl.Str("projectionMtx\x00")),
		1, false, &child2D.ProjectionMatrix[0],
	)

	gl.BindAttribLocation(child2D.ShaderProgram, 0, gl.Str("position\x00"))
	gl.BindAttribLocation(child2D.ShaderProgram, 1, gl.Str("tex\x00"))
}

func (child2D *Child2D) Update(mainCamera Camera) {
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.ShaderProgram, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(child2D.ShaderProgram, gl.Str("modelMtx\x00")),
		1, false, &child2D.ModelMatrix[0],
	)

	sX, sY := ScaleCoordinates(child2D.X, child2D.Y, float32(child2D.Config.ScreenWidth), float32(child2D.Config.ScreenHeight))
	child2D.ModelMatrix = mgl32.Translate3D(sX, sY, 0)
	child2D.ProjectionMatrix = mgl32.Ortho2D(-1, 1, -1, 1)
}

func (child2D *Child2D) AttachTexture(path string, coords []float32) error {
	if child2D.VertexArray == nil {
		return errors.New("Cannot attach texture without VertexArray")
	}
	if child2D.ShaderProgram == 0 {
		return errors.New("Cannot attach texture without shader program")
	}

	gl.BindVertexArray(child2D.VertexArray.id)
	gl.UseProgram(child2D.ShaderProgram)

	child2D.VertexArray.AddVertexAttribute(coords, 1, 2)

	texture, err := NewTexture(path)
	if err != nil {
		return err
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	loc1 := gl.GetUniformLocation(child2D.ShaderProgram, gl.Str("texture0\x00"))
	gl.Uniform1i(loc1, int32(0))
	CheckError("IGNORE")

	child2D.Texture = texture
	gl.BindVertexArray(0)
	return nil
}

func (child2D *Child2D) AttachTexturePrimitive(path string) {
	child2D.AttachTexture(path, GetPrimitiveCoords(child2D.Primitive))
}

func ScaleCoordinates(x, y, sw, sh float32) (float32, float32) {
	return 2*(x/float32(sw)) - 1, 2*(y/float32(sh)) - 1
}

func (child2D *Child2D) SetPosition(x, y float32) {
	child2D.X = x
	child2D.Y = y
}

func (child2D *Child2D) AttachVertexArray(vao *VertexArray, numVertices int32) {
	child2D.VertexArray = vao
	child2D.NumVertices = numVertices
}

func (child2D *Child2D) AttachPrimitive(p Primitive) {
	child2D.Primitive = p.id
	child2D.AttachVertexArray(p.vao, p.numVertices)
}

func (child2D *Child2D) AttachShader(s uint32) {
	child2D.ShaderProgram = s
}

func (child2D *Child2D) GetShaderProgram() uint32 {
	return child2D.ShaderProgram
}

func (child2D *Child2D) GetVertexArray() *VertexArray {
	return child2D.VertexArray
}

func (child2D *Child2D) GetNumVertices() int32 {
	return child2D.NumVertices
}

func (child2D *Child2D) GetTexture() uint32 {
	return child2D.Texture
}
