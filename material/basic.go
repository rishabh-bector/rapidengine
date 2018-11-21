package material

import "github.com/go-gl/gl/v4.1-core/gl"

type BasicMaterial struct {
	shader *ShaderProgram

	diffuseLevel float32

	color [3]float32

	diffuseMap      *uint32
	diffuseMapScale float32

	alphaLevel float32
	alphaMap   *uint32

	animationPlaying  string
	animationTextures map[string][]*uint32
	animationCurrent  int
	animationFrame    float64
	animationFPS      map[string]float64
	animationEnabled  bool
}

func NewBasicMaterial(shader *ShaderProgram) BasicMaterial {
	return BasicMaterial{
		shader:            shader,
		diffuseLevel:      0,
		color:             [3]float32{46, 46, 46},
		diffuseMapScale:   1,
		alphaLevel:        1,
		animationTextures: make(map[string][]*uint32),
		animationFPS:      make(map[string]float64),
		animationEnabled:  false,
	}
}

func (bm *BasicMaterial) Render(delta float64, darkness float32) {
	bm.UpdateAnimation(delta)
	bm.UpdateAttribArrays()

	gl.Uniform1f(bm.shader.GetUniform("diffuseLevel"), bm.diffuseLevel)

	gl.Uniform3fv(bm.shader.GetUniform("color"), 1, &bm.color[0])

	gl.Uniform1i(bm.shader.GetUniform("diffuseMap"), 0)
	gl.Uniform1f(bm.shader.GetUniform("diffuseMapScale"), bm.diffuseMapScale)

	gl.Uniform1f(bm.shader.GetUniform("alphaLevel"), bm.alphaLevel)
	gl.Uniform1i(bm.shader.GetUniform("alphaMap"), 1)

	gl.Uniform1f(bm.shader.GetUniform("darkness"), darkness)
}

func (bm *BasicMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
}

func (bm *BasicMaterial) UpdateAnimation(delta float64) {
	if bm.animationEnabled && bm.animationPlaying != "" {
		if bm.animationFrame > 1/bm.animationFPS[bm.animationPlaying] {
			if bm.animationCurrent < len(bm.animationTextures[bm.animationPlaying])-1 {
				bm.animationCurrent++
			} else {
				bm.animationCurrent = 0
			}
			bm.animationFrame = 0
			bm.diffuseMap = bm.animationTextures[bm.animationPlaying][bm.animationCurrent]
		} else {
			bm.animationFrame += delta
		}
	}
}
