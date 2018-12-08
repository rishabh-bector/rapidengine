package material

import "github.com/go-gl/gl/v4.1-core/gl"

type FoliageMaterial struct {
	shader *ShaderProgram

	// Foliage material
	DiffuseMap *uint32
	NormalMap  *uint32
	HeightMap  *uint32
}

func NewFoliageMaterial(shader *ShaderProgram) *FoliageMaterial {
	return &FoliageMaterial{
		shader: shader,
	}
}

func (fm *FoliageMaterial) Render(delta float64, darkness float32) {
	fm.UpdateAttribArrays()

	if fm.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *fm.DiffuseMap)
	}
	gl.Uniform1i(fm.shader.GetUniform("diffuseMap"), 0)

	if fm.NormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *fm.NormalMap)
	}
	gl.Uniform1i(fm.shader.GetUniform("normalMap"), 1)

	if fm.HeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *fm.HeightMap)
	}
	gl.Uniform1i(fm.shader.GetUniform("heightMap"), 2)
}

func (fm *FoliageMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	//gl.EnableVertexAttribArray(3)
	//gl.EnableVertexAttribArray(4)
}

func (fm *FoliageMaterial) GetShader() *ShaderProgram {
	return fm.shader
}
