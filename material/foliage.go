package material

import "github.com/go-gl/gl/v4.1-core/gl"

type FoliageMaterial struct {
	shader *ShaderProgram

	// Foliage material
	DiffuseMap *Texture
	NormalMap  *Texture
	HeightMap  *Texture
	OpacityMap *Texture

	// Terrain Info
	TerrainHeightMap    *Texture
	TerrainNormalMap    *Texture
	TerrainDisplacement float32

	TerrainWidth  float32
	TerrainLength float32

	FoliageDisplacement float32
	FoliageNoiseSeed    float32
	FoliageVariation    float32
}

func NewFoliageMaterial(shader *ShaderProgram) *FoliageMaterial {
	return &FoliageMaterial{
		shader:              shader,
		FoliageDisplacement: 1,
		FoliageNoiseSeed:    1,
		FoliageVariation:    0,
	}
}

func (fm *FoliageMaterial) Render(delta float64, darkness float32, totalTime float64) {
	if fm.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *fm.DiffuseMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("diffuseMap"), 0)

	if fm.NormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *fm.NormalMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("normalMap"), 1)

	if fm.HeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *fm.HeightMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("heightMap"), 2)

	if fm.OpacityMap != nil {
		gl.ActiveTexture(gl.TEXTURE5)
		gl.BindTexture(gl.TEXTURE_2D, *fm.OpacityMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("opacityMap"), 5)

	if fm.TerrainHeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *fm.TerrainHeightMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("terrainHeightMap"), 3)

	if fm.TerrainNormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_2D, *fm.TerrainNormalMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("terrainNormalMap"), 4)

	gl.Uniform1f(fm.shader.GetUniform("terrainDisplacement"), fm.TerrainDisplacement)
	gl.Uniform1f(fm.shader.GetUniform("terrainWidth"), fm.TerrainWidth)
	gl.Uniform1f(fm.shader.GetUniform("terrainLength"), fm.TerrainLength)

	gl.Uniform1f(fm.shader.GetUniform("foliageDisplacement"), fm.FoliageDisplacement)
	gl.Uniform1f(fm.shader.GetUniform("foliageNoiseSeed"), fm.FoliageNoiseSeed)
	gl.Uniform1f(fm.shader.GetUniform("foliageVariation"), fm.FoliageVariation)

	gl.Uniform1f(fm.shader.GetUniform("totalTime"), float32(totalTime))
}

func (fm *FoliageMaterial) GetShader() *ShaderProgram {
	return fm.shader
}
