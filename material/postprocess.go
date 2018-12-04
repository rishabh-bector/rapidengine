package material

import "github.com/go-gl/gl/v4.1-core/gl"

type PostProcessMaterial struct {
	shader *ShaderProgram

	ScreenMap *uint32
}

func NewPostProcessMaterial(shader *ShaderProgram, screenMap *uint32) *PostProcessMaterial {
	return &PostProcessMaterial{
		shader:    shader,
		ScreenMap: screenMap,
	}
}

func (pm *PostProcessMaterial) Render(delta float64, darkness float32) {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *pm.ScreenMap)

	gl.Uniform1i(pm.shader.GetUniform("screen"), 0)
}

func (pm *PostProcessMaterial) GetShader() *ShaderProgram {
	return pm.shader
}
