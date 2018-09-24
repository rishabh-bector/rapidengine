package rapidengine

import (
	"fmt"
	"rapidengine/camera"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SkyBox struct {
	shader   uint32
	material Material

	vao *VertexArray

	projectionMatrix mgl32.Mat4
	modelMatrix      mgl32.Mat4
}

func NewSkyBox(path string, engine *Engine) {
	gl.UseProgram(engine.ShaderControl.GetShader("skybox"))
	gl.BindAttribLocation(engine.ShaderControl.GetShader("skybox"), 0, gl.Str("position\x00"))

	engine.TextureControl.NewCubeMap(
		fmt.Sprintf("../rapidengine/skybox/%s/%s_LF.png", path, path),
		fmt.Sprintf("../rapidengine/skybox/%s/%s_RT.png", path, path),
		fmt.Sprintf("../rapidengine/skybox/%s/%s_UP.png", path, path),
		fmt.Sprintf("../rapidengine/skybox/%s/%s_DN.png", path, path),
		fmt.Sprintf("../rapidengine/skybox/%s/%s_FR.png", path, path),
		fmt.Sprintf("../rapidengine/skybox/%s/%s_BK.png", path, path),
		"skybox")

	material := NewMaterial(engine.ShaderControl.GetShader("skybox"))
	material.BecomeCubemap(engine.TextureControl.GetTexture("skybox"))

	indices := []uint32{}
	for i := 0; i < len(skyBoxVertices); i++ {
		indices = append(indices, uint32(i))
	}

	vao := NewVertexArray(skyBoxVertices, indices)
	vao.AddVertexAttribute(CubeTextures, 1, 2)

	engine.Renderer.SkyBoxEnabled = true
	engine.Renderer.SkyBox = &SkyBox{
		shader:   engine.ShaderControl.GetShader("skybox"),
		material: material,
		vao:      vao,
		projectionMatrix: mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(engine.Config.ScreenWidth)/float32(engine.Config.ScreenHeight),
			0.1, 100,
		),
		modelMatrix: mgl32.Ident4(),
	}
}

func (skyBox *SkyBox) Render(mainCamera camera.Camera) {
	gl.DepthMask(false)
	gl.UseProgram(skyBox.shader)
	gl.BindVertexArray(skyBox.vao.id)

	skyBox.material.Render()

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
