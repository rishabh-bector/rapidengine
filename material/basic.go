package material

import (
	"rapidengine/state"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type BasicMaterial struct {
	Shader *ShaderProgram

	DiffuseLevel float32

	Hue [4]float32

	DiffuseMap      *uint32
	DiffuseMapScale float32

	AlphaMapLevel float32
	AlphaMap      *uint32

	Flipped int

	ScatterLevel float32

	animationPlaying  string
	animationTextures map[string][]*uint32
	animationCurrent  int
	animationFrame    float64
	animationFPS      map[string]float64
	animationEnabled  bool

	animationPlayingOnce  bool
	animationOnceCallback func()
}

func NewBasicMaterial(Shader *ShaderProgram) *BasicMaterial {
	return &BasicMaterial{
		Shader:            Shader,
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

func (bm *BasicMaterial) Render(delta float64, darkness float32, totalTime float64) {
	bm.UpdateAnimation(delta)

	if bm.DiffuseMap != nil && state.BoundTexture0 != *bm.DiffuseMap {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *bm.DiffuseMap)
		state.BoundTexture0 = *bm.DiffuseMap
	}

	if bm.AlphaMap != nil && state.BoundTexture1 != *bm.AlphaMap {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, *bm.AlphaMap)
		state.BoundTexture1 = *bm.AlphaMap
	}

	gl.Uniform1f(bm.Shader.GetUniform("diffuseLevel"), bm.DiffuseLevel)

	gl.Uniform4fv(bm.Shader.GetUniform("hue"), 1, &bm.Hue[0])

	gl.Uniform1i(bm.Shader.GetUniform("diffuseMap"), 0)
	gl.Uniform1f(bm.Shader.GetUniform("scale"), bm.DiffuseMapScale)

	gl.Uniform1f(bm.Shader.GetUniform("alphaMapLevel"), bm.AlphaMapLevel)
	gl.Uniform1i(bm.Shader.GetUniform("alphaMap"), 1)

	gl.Uniform1f(bm.Shader.GetUniform("darkness"), darkness)

	gl.Uniform1f(bm.Shader.GetUniform("scatterLevel"), bm.ScatterLevel)

	gl.Uniform1i(bm.Shader.GetUniform("flipped"), int32(bm.Flipped))
}

func (bm *BasicMaterial) GetShader() *ShaderProgram {
	return bm.Shader
}

func (bm *BasicMaterial) UpdateAnimation(delta float64) {
	if bm.animationEnabled && bm.animationPlaying != "" {
		if bm.animationFrame > 1/bm.animationFPS[bm.animationPlaying] {
			if bm.animationCurrent < len(bm.animationTextures[bm.animationPlaying])-1 {
				bm.animationCurrent++
			} else {
				if bm.animationPlayingOnce {
					bm.animationPlaying = ""
					bm.animationPlayingOnce = false
					if bm.animationOnceCallback != nil {
						bm.animationOnceCallback()
						bm.animationOnceCallback = nil
					}
					return
				}
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
	bm.animationCurrent = 0
	bm.animationFrame = 0
	bm.animationPlayingOnce = false
	bm.DiffuseMap = bm.animationTextures[bm.animationPlaying][bm.animationCurrent]
}

func (bm *BasicMaterial) PlayAnimationOnce(anim string) {
	bm.PlayAnimation(anim)
	bm.animationPlayingOnce = true
}

func (bm *BasicMaterial) PlayAnimationOnceCallback(anim string, callback func()) {
	bm.PlayAnimation(anim)
	bm.animationPlayingOnce = true
	bm.animationOnceCallback = callback
}
