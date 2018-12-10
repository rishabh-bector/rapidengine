package material

import "github.com/go-gl/gl/v4.1-core/gl"

type TerrainMaterial struct {
	shader *ShaderProgram

	// Terrain material
	DiffuseMap *uint32
	NormalMap  *uint32
	HeightMap  *uint32

	// Terrain data
	TerrainHeightMap    *uint32
	TerrainNormalMap    *uint32
	TerrainDisplacement float32

	Displacement float32
	Scale        float32
}

func NewTerrainMaterial(shader *ShaderProgram) *TerrainMaterial {
	return &TerrainMaterial{
		shader: shader,
		Scale:  1,
	}
}

func (tm *TerrainMaterial) Render(delta float64, darkness float32, totalTime float64) {
	tm.UpdateAttribArrays()

	if tm.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *tm.DiffuseMap)
	}
	gl.Uniform1i(tm.shader.GetUniform("diffuseMap"), 0)

	if tm.NormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *tm.NormalMap)
	}
	gl.Uniform1i(tm.shader.GetUniform("normalMap"), 1)

	if tm.HeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *tm.HeightMap)
	}
	gl.Uniform1i(tm.shader.GetUniform("heightMap"), 2)

	if tm.TerrainHeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *tm.TerrainHeightMap)
	}
	gl.Uniform1i(tm.shader.GetUniform("terrainHeightMap"), 3)

	if tm.TerrainNormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_2D, *tm.TerrainNormalMap)
	}
	gl.Uniform1i(tm.shader.GetUniform("terrainNormalMap"), 4)

	gl.Uniform1f(tm.shader.GetUniform("terrainDisplacement"), tm.TerrainDisplacement)

	gl.Uniform1f(tm.shader.GetUniform("displacement"), tm.Displacement)
	gl.Uniform1f(tm.shader.GetUniform("scale"), tm.Scale)
}

func (tm *TerrainMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	//gl.EnableVertexAttribArray(3)
	//gl.EnableVertexAttribArray(4)
}

func (tm *TerrainMaterial) GetShader() *ShaderProgram {
	return tm.shader
}
