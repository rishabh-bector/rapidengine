package material

import "github.com/go-gl/gl/v4.1-core/gl"

type StandardMaterial struct {
	shader *ShaderProgram

	diffuseMap  *uint32
	normalMap   *uint32
	heightMap   *uint32
	specularMap *uint32

	Hue          [4]float32
	DiffuseLevel float32

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
	}
}

func (sm *StandardMaterial) AttachDiffuseMap(dm *uint32) {
	sm.diffuseMap = dm
}

func (sm *StandardMaterial) AttachNormalMap(nm *uint32) {
	sm.normalMap = nm
}

func (sm *StandardMaterial) AttachHeightMap(hm *uint32) {
	sm.heightMap = hm
}

func (sm *StandardMaterial) AttachSpecularMap(hm *uint32) {
	sm.specularMap = hm
}

func (sm *StandardMaterial) Render(delta float64, darkness float32, totalTime float64) {

	if sm.diffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *sm.diffuseMap)
	}
	gl.Uniform1i(sm.shader.GetUniform("diffuseMap"), 0)

	if sm.normalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *sm.normalMap)
	}
	gl.Uniform1i(sm.shader.GetUniform("normalMap"), 1)

	if sm.heightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *sm.heightMap)
	}
	gl.Uniform1i(sm.shader.GetUniform("heightMap"), 2)

	if sm.specularMap != nil {
		gl.ActiveTexture(gl.TEXTURE3)
		gl.BindTexture(gl.TEXTURE_2D, *sm.specularMap)
	}
	gl.Uniform1i(sm.shader.GetUniform("specularMap"), 3)

	gl.Uniform4fv(sm.shader.GetUniform("hue"), 1, &sm.Hue[0])
	gl.Uniform1f(sm.shader.GetUniform("diffuseLevel"), sm.DiffuseLevel)

	gl.Uniform1f(sm.shader.GetUniform("displacement"), sm.Displacement)
	gl.Uniform1f(sm.shader.GetUniform("scale"), sm.Scale)

	gl.Uniform1f(sm.shader.GetUniform("reflectivity"), sm.Reflectivity)
	gl.Uniform1f(sm.shader.GetUniform("refractivity"), sm.Refractivity)
	gl.Uniform1f(sm.shader.GetUniform("refractLevel"), sm.RefractLevel)
}

func (sm *StandardMaterial) GetShader() *ShaderProgram {
	return sm.shader
}
