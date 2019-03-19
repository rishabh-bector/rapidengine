package material

import (
	"rapidengine/state"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type PostProcessMaterial struct {
	shader *ShaderProgram

	FboWidth  float32
	FboHeight float32

	ScreenMap *uint32
}

func NewPostProcessMaterial(shader *ShaderProgram, screenMap *uint32) *PostProcessMaterial {
	return &PostProcessMaterial{
		shader:    shader,
		ScreenMap: screenMap,
	}
}

func (pm *PostProcessMaterial) Render(delta float64, darkness float32, totalTime float64) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *pm.ScreenMap)
	state.BoundTexture0 = *pm.ScreenMap

	gl.Uniform1i(pm.shader.GetUniform("screen"), 0)

	gl.Uniform1f(pm.shader.GetUniform("fboWidth"), pm.FboWidth)
	gl.Uniform1f(pm.shader.GetUniform("fboHeight"), pm.FboHeight)
}

func (pm *PostProcessMaterial) GetShader() *ShaderProgram {
	return pm.shader
}

func (pm *PostProcessMaterial) AttachShader(shader *ShaderProgram) {
	pm.shader = shader
}
