package material

import "github.com/go-gl/gl/v4.3-core/gl"

type TerrainMaterial struct {
	shader *ShaderProgram

	// Terrain material
	DiffuseMap *Texture
	NormalMap  *Texture
	HeightMap  *Texture

	// Terrain data
	TerrainHeightMap    *Texture
	TerrainNormalMap    *Texture
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
	if tm.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *tm.DiffuseMap.Addr)
	}
	gl.Uniform1i(tm.shader.GetUniform("diffuseMap"), 0)

	if tm.NormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *tm.NormalMap.Addr)
	}
	gl.Uniform1i(tm.shader.GetUniform("normalMap"), 1)

	if tm.HeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *tm.HeightMap.Addr)
	}
	gl.Uniform1i(tm.shader.GetUniform("heightMap"), 2)

	if tm.TerrainHeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *tm.TerrainHeightMap.Addr)
	}
	gl.Uniform1i(tm.shader.GetUniform("terrainHeightMap"), 3)

	if tm.TerrainNormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_2D, *tm.TerrainNormalMap.Addr)
	}
	gl.Uniform1i(tm.shader.GetUniform("terrainNormalMap"), 4)

	gl.Uniform1f(tm.shader.GetUniform("terrainDisplacement"), tm.TerrainDisplacement)

	gl.Uniform1f(tm.shader.GetUniform("displacement"), tm.Displacement)
	gl.Uniform1f(tm.shader.GetUniform("scale"), tm.Scale)
}

func (tm *TerrainMaterial) GetShader() *ShaderProgram {
	return tm.shader
}
