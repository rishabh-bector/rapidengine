package material

import "github.com/go-gl/gl/v4.3-core/gl"

type WaterMaterial struct {
	shader *ShaderProgram

	diffuseMap *uint32
	normalMap  *uint32
	heightMap  *uint32

	displacement float32
	Scale        float32
}

func NewWaterMaterial(shader *ShaderProgram) *WaterMaterial {
	return &WaterMaterial{
		shader: shader,
		Scale:  1,
	}
}

func (wm *WaterMaterial) AttachDiffuseMap(dm *uint32) {
	wm.diffuseMap = dm
}

func (wm *WaterMaterial) AttachNormalMap(nm *uint32) {
	wm.normalMap = nm
}

func (wm *WaterMaterial) AttachHeightMap(hm *uint32) {
	wm.heightMap = hm
}

func (wm *WaterMaterial) SetDisplacement(d float32) {
	wm.displacement = d
}

func (wm *WaterMaterial) Render(delta float64, darkness float32, totalTime float64) {
	wm.UpdateAttribArrays()

	if wm.diffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *wm.diffuseMap)
	}

	gl.Uniform1i(wm.shader.GetUniform("diffuseMap"), 0)

	if wm.normalMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *wm.normalMap)
	}

	gl.Uniform1i(wm.shader.GetUniform("normalMap"), 1)

	if wm.heightMap != nil {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, *wm.heightMap)
	}

	gl.Uniform1i(wm.shader.GetUniform("heightMap"), 2)

	gl.Uniform1f(wm.shader.GetUniform("displacement"), wm.displacement)
	gl.Uniform1f(wm.shader.GetUniform("scale"), wm.Scale)

	gl.Uniform1f(wm.shader.GetUniform("totalTime"), float32(totalTime))
}

func (wm *WaterMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	gl.EnableVertexAttribArray(4)
}

func (wm *WaterMaterial) GetShader() *ShaderProgram {
	return wm.shader
}
