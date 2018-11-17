package terrain

import (
	"rapidengine/camera"
	"rapidengine/geometry"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SkyBox struct {
	shader   *material.ShaderProgram
	material material.Material

	vao *geometry.VertexArray

	projectionMatrix mgl32.Mat4
	modelMatrix      mgl32.Mat4
}

func NewSkyBox(
	shader *material.ShaderProgram,
	mat material.Material,
	vao *geometry.VertexArray,
	projectionMatrix,
	modelMatrix mgl32.Mat4,
) *SkyBox {
	return &SkyBox{
		shader:           shader,
		material:         mat,
		vao:              vao,
		projectionMatrix: projectionMatrix,
		modelMatrix:      modelMatrix,
	}
}

func (skyBox *SkyBox) Render(mainCamera camera.Camera) {
	gl.DepthMask(false)
	skyBox.shader.Bind()
	gl.BindVertexArray(skyBox.vao.GetID())

	skyBox.material.Render(0, 1)

	x, y, z := mainCamera.GetPosition()
	skyBox.modelMatrix = mgl32.Translate3D(x, y, z)

	gl.UniformMatrix4fv(
		skyBox.shader.GetUniform("modelMtx"),
		1, false, &skyBox.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		skyBox.shader.GetUniform("viewMtx"),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		skyBox.shader.GetUniform("projectionMtx"),
		1, false, &skyBox.projectionMatrix[0],
	)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.DrawElements(gl.TRIANGLES, 108, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.DepthMask(true)
}

var SkyBoxVertices = []float32{

	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	// R1
	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}
