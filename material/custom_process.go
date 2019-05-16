package material

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CustomProcessMaterial struct {
	shader    *ShaderProgram
	ScreenMap *uint32

	// Custom func for user uniforms
	RenderFunc func(delta float64, darkness float32, totalTime float64)

	FboWidth  float32
	FboHeight float32
}

func NewCustomProcessMaterial(shader *ShaderProgram) *CustomProcessMaterial {
	return &CustomProcessMaterial{
		shader: shader,
	}
}

func (sm *CustomProcessMaterial) Render(delta float64, darkness float32, totalTime float64) {
	gl.Uniform1f(sm.shader.GetUniform("fboWidth"), sm.FboWidth)
	gl.Uniform1f(sm.shader.GetUniform("fboHeight"), sm.FboHeight)

	sm.RenderFunc(delta, darkness, totalTime)
}

func (sm *CustomProcessMaterial) BindCustomInput(index uint32, texture uint32, uniform string) {
	gl.ActiveTexture(gl.TEXTURE0 + index)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.Uniform1i(sm.shader.GetUniform(uniform), int32(index))
}

func (sm *CustomProcessMaterial) GetShader() *ShaderProgram {
	return sm.shader
}
