package material

import "github.com/go-gl/gl/v4.1-core/gl"

type PBRMaterial struct {
	shader *ShaderProgram

	AlbedoMap           *uint32
	NormalMap           *uint32
	HeightMap           *uint32
	MetallicMap         *uint32
	RoughnessMap        *uint32
	AmbientOcclusionMap *uint32

	NormalScalar           float32
	MetallicScalar         float32
	RoughnessScalar        float32
	AmbientOcclusionScalar float32

	VertexDisplacement   float32
	ParallaxDisplacement float32

	Scale float32

	Reflectivity float32
	Refractivity float32
	RefractLevel float32
}

func NewPBRMaterial(shader *ShaderProgram) *PBRMaterial {
	return &PBRMaterial{
		shader: shader,
		Scale:  1,

		MetallicScalar:         -0.5,
		RoughnessScalar:        -0.5,
		AmbientOcclusionScalar: 1,
	}
}

func (pm *PBRMaterial) Render(delta float64, darkness float32, totalTime float64) {

	if pm.AlbedoMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *pm.AlbedoMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("albedoMap"), 0)

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

	if pm.MetallicMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *pm.MetallicMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("metallicMap"), 3)

	if pm.RoughnessMap != nil {
		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_2D, *pm.RoughnessMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("roughnessMap"), 4)

	if pm.AmbientOcclusionMap != nil {
		gl.ActiveTexture(gl.TEXTURE5)
		gl.BindTexture(gl.TEXTURE_2D, *pm.AmbientOcclusionMap)
	}
	gl.Uniform1i(pm.shader.GetUniform("aoMap"), 5)

	gl.Uniform1f(pm.shader.GetUniform("normalScalar"), pm.NormalScalar)
	gl.Uniform1f(pm.shader.GetUniform("metallicScalar"), pm.MetallicScalar)
	gl.Uniform1f(pm.shader.GetUniform("roughnessScalar"), pm.RoughnessScalar)
	gl.Uniform1f(pm.shader.GetUniform("aoScalar"), pm.AmbientOcclusionScalar)

	gl.Uniform1f(pm.shader.GetUniform("scale"), pm.Scale)
	gl.Uniform1f(pm.shader.GetUniform("vertexDisplacement"), pm.VertexDisplacement)
	gl.Uniform1f(pm.shader.GetUniform("parallaxDisplacement"), pm.ParallaxDisplacement)

	gl.Uniform1f(pm.shader.GetUniform("reflectivity"), pm.Reflectivity)
	gl.Uniform1f(pm.shader.GetUniform("refractivity"), pm.Refractivity)
	gl.Uniform1f(pm.shader.GetUniform("refractLevel"), pm.RefractLevel)
}

func (pm *PBRMaterial) GetShader() *ShaderProgram {
	return pm.shader
}
