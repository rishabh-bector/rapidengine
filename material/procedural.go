package material

import "github.com/go-gl/gl/v4.1-core/gl"

type ProceduralMaterial struct {
	shader *ShaderProgram

	// Terrain material
	DiffuseMap *uint32
	NormalMap  *uint32
	HeightMap  *uint32

	Displacement float32
	Scale        float32
}

func NewProceduralMaterial(shader *ShaderProgram) *ProceduralMaterial {
	return &ProceduralMaterial{
		shader: shader,
		Scale:  1,
	}
}

func (pm *ProceduralMaterial) Render(delta float64, darkness float32, totalTime float64) {
	pm.UpdateAttribArrays()

	if pm.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *pm.DiffuseMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("diffuseMap"), 0)

	if pm.NormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *pm.NormalMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("normalMap"), 1)

	if pm.HeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *pm.HeightMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("heightMap"), 2)

	gl.Uniform1f(pm.shader.GetUniform("displacement"), pm.Displacement)
	gl.Uniform1f(pm.shader.GetUniform("scale"), pm.Scale)
}

func (pm *ProceduralMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	//gl.EnableVertexAttribArray(3)
	//gl.EnableVertexAttribArray(4)
}

func (pm *ProceduralMaterial) GetShader() *ShaderProgram {
	return pm.shader
}
