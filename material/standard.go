package material

import "github.com/go-gl/gl/v4.3-core/gl"

type StandardMaterial struct {
	shader *ShaderProgram

	diffuseMap  *Texture
	normalMap   *Texture
	heightMap   *Texture
	specularMap *Texture

	DiffuseLevel  float32
	NormalLevel   float32
	SpecularLevel float32
	HeightLevel   float32

	Hue [4]float32

	Displacement float32
	Scale        float32

	Reflectivity float32
	Refractivity float32
	RefractLevel float32
}

func NewStandardMaterial(shader *ShaderProgram) *StandardMaterial {
	return &StandardMaterial{
		shader: shader,
		Scale:  1,
		Hue:    [4]float32{100, 100, 100, 255},

		SpecularLevel: 1,
	}
}

func (sm *StandardMaterial) AttachDiffuseMap(dm *Texture) {
	sm.diffuseMap = dm
}

func (sm *StandardMaterial) AttachNormalMap(nm *Texture) {
	sm.normalMap = nm
}

func (sm *StandardMaterial) AttachHeightMap(hm *Texture) {
	sm.heightMap = hm
}

func (sm *StandardMaterial) AttachSpecularMap(hm *Texture) {
	sm.specularMap = hm
}

func (sm *StandardMaterial) Render(delta float64, darkness float32, totalTime float64) {

	if sm.diffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *sm.diffuseMap.Addr)
	}
	gl.Uniform1i(sm.shader.GetUniform("diffuseMap"), 0)

	if sm.normalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *sm.normalMap.Addr)
	}
	gl.Uniform1i(sm.shader.GetUniform("normalMap"), 1)

	if sm.heightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *sm.heightMap.Addr)
	}
	gl.Uniform1i(sm.shader.GetUniform("heightMap"), 2)

	if sm.specularMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *sm.specularMap.Addr)
	}
	gl.Uniform1i(sm.shader.GetUniform("specularMap"), 3)

	gl.Uniform1f(sm.shader.GetUniform("diffuseLevel"), sm.DiffuseLevel)
	gl.Uniform1f(sm.shader.GetUniform("normalLevel"), sm.NormalLevel)
	gl.Uniform1f(sm.shader.GetUniform("specularLevel"), sm.SpecularLevel)
	gl.Uniform1f(sm.shader.GetUniform("heightLevel"), sm.HeightLevel)

	gl.Uniform4fv(sm.shader.GetUniform("hue"), 1, &sm.Hue[0])

	gl.Uniform1f(sm.shader.GetUniform("displacement"), sm.Displacement)
	gl.Uniform1f(sm.shader.GetUniform("scale"), sm.Scale)

	gl.Uniform1f(sm.shader.GetUniform("reflectivity"), sm.Reflectivity)
	gl.Uniform1f(sm.shader.GetUniform("refractivity"), sm.Refractivity)
	gl.Uniform1f(sm.shader.GetUniform("refractLevel"), sm.RefractLevel)
}

func (sm *StandardMaterial) GetShader() *ShaderProgram {
	return sm.shader
}
