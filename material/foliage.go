package material

import "github.com/go-gl/gl/v4.3-core/gl"

type FoliageMaterial struct {
	shader *ShaderProgram

	// Standard material
	DiffuseMap  *Texture
	NormalMap   *Texture
	HeightMap   *Texture
	SpecularMap *Texture
	OpacityMap  *Texture

	DiffuseLevel  float32
	NormalLevel   float32
	HeightLevel   float32
	SpecularLevel float32

	Displacement float32
	Scale        float32
	Hue          [4]float32

	Reflectivity float32
	Refractivity float32
	RefractLevel float32

	// Terrain Info
	TerrainHeightMap    *Texture
	TerrainNormalMap    *Texture
	TerrainDisplacement float32

	TerrainWidth  float32
	TerrainLength float32

	// Foliage Config
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

		Scale:         1.0,
		Hue:           [4]float32{100, 100, 100, 255},
		SpecularLevel: 1,
	}
}

func (fm *FoliageMaterial) Render(delta float64, darkness float32, totalTime float64) {

	//   --------------------------------------------------
	//   Standard Material
	//   --------------------------------------------------

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

	if fm.SpecularMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *fm.SpecularMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("specularMap"), 3)

	gl.Uniform1f(fm.shader.GetUniform("diffuseLevel"), fm.DiffuseLevel)
	gl.Uniform1f(fm.shader.GetUniform("normalLevel"), fm.NormalLevel)
	gl.Uniform1f(fm.shader.GetUniform("specularLevel"), fm.SpecularLevel)
	gl.Uniform1f(fm.shader.GetUniform("heightLevel"), fm.HeightLevel)

	gl.Uniform4fv(fm.shader.GetUniform("hue"), 1, &fm.Hue[0])

	gl.Uniform1f(fm.shader.GetUniform("displacement"), fm.Displacement)
	gl.Uniform1f(fm.shader.GetUniform("scale"), fm.Scale)

	gl.Uniform1f(fm.shader.GetUniform("reflectivity"), fm.Reflectivity)
	gl.Uniform1f(fm.shader.GetUniform("refractivity"), fm.Refractivity)
	gl.Uniform1f(fm.shader.GetUniform("refractLevel"), fm.RefractLevel)

	//   --------------------------------------------------
	//   Foliage Material
	//   --------------------------------------------------

	if fm.OpacityMap != nil {
		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_2D, *fm.OpacityMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("opacityMap"), 4)

	if fm.TerrainHeightMap != nil {
		gl.ActiveTexture(gl.TEXTURE5)
		gl.BindTexture(gl.TEXTURE_2D, *fm.TerrainHeightMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("terrainHeightMap"), 5)

	if fm.TerrainNormalMap != nil {
		gl.ActiveTexture(gl.TEXTURE6)
		gl.BindTexture(gl.TEXTURE_2D, *fm.TerrainNormalMap.Addr)
	}
	gl.Uniform1i(fm.shader.GetUniform("terrainNormalMap"), 6)

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
