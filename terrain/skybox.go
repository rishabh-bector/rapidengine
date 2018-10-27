package terrain

import (
	"fmt"
	"rapidengine/camera"
	"rapidengine/configuration"
	"rapidengine/control"
	"rapidengine/geometry"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SkyBox struct {
	shader   uint32
	material material.Material

	vao *geometry.VertexArray

	projectionMatrix mgl32.Mat4
	modelMatrix      mgl32.Mat4
}

func NewSkyBox(path string, shaderControl *control.ShaderControl, textureControl *control.TextureControl, config *configuration.EngineConfig) *SkyBox {
	gl.UseProgram(shaderControl.GetShader("skybox"))
	gl.BindAttribLocation(shaderControl.GetShader("skybox"), 0, gl.Str("position\x00"))

	textureControl.NewCubeMap(
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_LF.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_RT.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_UP.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_DN.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_FR.png", path, path),
		fmt.Sprintf("../rapidengine/assets/skybox/%s/%s_BK.png", path, path),
		"skybox")

	material := material.NewMaterial(shaderControl.GetShader("skybox"), config)
	material.BecomeCubemap(textureControl.GetTexture("skybox"))

	indices := []uint32{}
	for i := 0; i < len(skyBoxVertices); i++ {
		indices = append(indices, uint32(i))
	}

	vao := geometry.NewVertexArray(skyBoxVertices, indices)
	vao.AddVertexAttribute(geometry.CubeTextures, 1, 2)

	return &SkyBox{
		shader:   shaderControl.GetShader("skybox"),
		material: material,
		vao:      vao,
		projectionMatrix: mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(config.ScreenWidth)/float32(config.ScreenHeight),
			0.1, 100,
		),
		modelMatrix: mgl32.Ident4(),
	}
}

func (skyBox *SkyBox) Render(mainCamera camera.Camera) {
	gl.DepthMask(false)
	gl.UseProgram(skyBox.shader)
	gl.BindVertexArray(skyBox.vao.GetID())

	skyBox.material.Render(0, 1)

	x, y, z := mainCamera.GetPosition()
	skyBox.modelMatrix = mgl32.Translate3D(x, y, z)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(skyBox.shader, gl.Str("modelMtx\x00")),
		1, false, &skyBox.modelMatrix[0],
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(skyBox.shader, gl.Str("viewMtx\x00")),
		1, false, mainCamera.GetFirstViewIndex(),
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(skyBox.shader, gl.Str("projectionMtx\x00")),
		1, false, &skyBox.projectionMatrix[0],
	)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.DrawElements(gl.TRIANGLES, 108, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.DepthMask(true)
}

var skyBoxVertices = []float32{

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
