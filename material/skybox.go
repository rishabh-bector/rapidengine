package material

import "github.com/go-gl/gl/v4.1-core/gl"

type CubemapMaterial struct {
	shader *ShaderProgram

	CubeDiffuseMap *uint32
}

func NewCubemapMaterial(shader *ShaderProgram) *CubemapMaterial {
	return &CubemapMaterial{
		shader: shader,
	}
}

func (cm *CubemapMaterial) Render(delta float64, darkness float32, totalTime float64) {
	if cm.CubeDiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *cm.CubeDiffuseMap)
	}
}

func (cm *CubemapMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
}

func (cm *CubemapMaterial) GetShader() *ShaderProgram {
	return cm.shader
}
