package material

import "github.com/go-gl/gl/v4.1-core/gl"

type BasicMaterial struct {
	shader *ShaderProgram

	DiffuseLevel float32

	Hue [4]float32

	DiffuseMap      *uint32
	DiffuseMapScale float32

	AlphaMapLevel float32
	AlphaMap      *uint32

	animationPlaying  string
	animationTextures map[string][]*uint32
	animationCurrent  int
	animationFrame    float64
	animationFPS      map[string]float64
	animationEnabled  bool
}

func NewBasicMaterial(shader *ShaderProgram) *BasicMaterial {
	return &BasicMaterial{
		shader:            shader,
		DiffuseLevel:      0,
		Hue:               [4]float32{200, 200, 200, 255},
		DiffuseMapScale:   1,
		AlphaMapLevel:     0,
		animationTextures: make(map[string][]*uint32),
		animationFPS:      make(map[string]float64),
		animationPlaying:  "",
		animationEnabled:  false,
	}
}

func (bm *BasicMaterial) Render(delta float64, darkness float32) {
	bm.UpdateAnimation(delta)
	bm.UpdateAttribArrays()

	if bm.DiffuseMap != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *bm.DiffuseMap)
	}

	if bm.AlphaMap != nil {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *bm.AlphaMap)
	}

	gl.Uniform1f(bm.shader.GetUniform("diffuseLevel"), bm.DiffuseLevel)

	gl.Uniform4fv(bm.shader.GetUniform("hue"), 1, &bm.Hue[0])

	gl.Uniform1i(bm.shader.GetUniform("diffuseMap"), 0)
	gl.Uniform1f(bm.shader.GetUniform("diffuseMapScale"), bm.DiffuseMapScale)

	gl.Uniform1f(bm.shader.GetUniform("alphaMapLevel"), bm.AlphaMapLevel)
	gl.Uniform1i(bm.shader.GetUniform("alphaMap"), 1)

	gl.Uniform1f(bm.shader.GetUniform("darkness"), darkness)
}

func (bm *BasicMaterial) UpdateAttribArrays() {
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
}

func (bm *BasicMaterial) GetShader() *ShaderProgram {
	return bm.shader
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
			bm.DiffuseMap = bm.animationTextures[bm.animationPlaying][bm.animationCurrent]
		} else {
			bm.animationFrame += delta
		}
	}
}

func (bm *BasicMaterial) EnableAnimation() {
	bm.animationEnabled = true
}

func (bm *BasicMaterial) AddFrame(frame *uint32, anim string) {
	bm.animationTextures[anim] = append(bm.animationTextures[anim], frame)
}

func (bm *BasicMaterial) SetAnimationFPS(anim string, s float64) {
	bm.animationFPS[anim] = s
}

func (bm *BasicMaterial) PlayAnimation(anim string) {
	bm.animationPlaying = anim
}
