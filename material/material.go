package material

import (
	"rapidengine/configuration"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const SHADER_COLOR = "SHADER_COLOR"
const SHADER_TEXTURE = "SHADER_TEXTURE"
const SHADER_CUBEMAP = "SHADER_CUBEMAP"

var SHADER_COLOR_UNI = []float32{1, 0, 0}
var SHADER_TEXTURE_UNI = []float32{0, 1, 0}
var SHADER_CUBEMAP_UNI = []float32{0, 0, 1}

type Material struct {
	shader     *ShaderProgram
	shaderType string

	texture      *uint32
	textureScale float32

	transparencyEnabled bool
	transparencyTexture *uint32

	transparency float32

	color [3]float32
	shine float32

	animationPlaying  string
	animationTextures map[string][]*uint32
	animationCurrent  int
	animationFrame    float64
	animationFPS      map[string]float64
	animationEnabled  bool

	off bool
}

func NewMaterial(program *ShaderProgram, config *configuration.EngineConfig) Material {
	m := Material{
		shader:              program,
		color:               [3]float32{1, 1, 1},
		shine:               0.8,
		textureScale:        1,
		transparencyEnabled: false,
		animationEnabled:    false,
		animationCurrent:    0,
		animationTextures:   make(map[string][]*uint32),
		animationFPS:        make(map[string]float64),
		transparency:        1,
		off:                 false,
	}
	if config.SingleMaterial {
		m.off = true
	}
	return m
}

func (material *Material) PreRender() {
	switch material.shaderType {
	case SHADER_COLOR:
	case SHADER_TEXTURE:
	}
}

func (material *Material) Render(delta float64, darkness float32) {
	if material.off {
		return
	}

	if material.animationEnabled && material.animationPlaying != "" {
		if material.animationFrame > 1/material.animationFPS[material.animationPlaying] {
			if material.animationCurrent < len(material.animationTextures[material.animationPlaying])-1 {
				material.animationCurrent++
			} else {
				material.animationCurrent = 0
			}
			material.animationFrame = 0
			material.texture = material.animationTextures[material.animationPlaying][material.animationCurrent]
		} else {
			material.animationFrame += delta
		}
	}

	switch material.shaderType {

	case SHADER_COLOR:
		gl.Uniform3fv(material.shader.GetUniform("color"), 1, &material.color[0])
		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)
		gl.EnableVertexAttribArray(2)

		gl.Uniform3fv(material.shader.GetUniform("materialType"), 1, &SHADER_COLOR_UNI[0])
		gl.Uniform1i(material.shader.GetUniform("diffuseMap"), 0)
		gl.Uniform1i(material.shader.GetUniform("cubeDiffuseMap"), 1)
		gl.Uniform1f(material.shader.GetUniform("shine"), material.shine)
		gl.Uniform1f(material.shader.GetUniform("textureScale"), material.textureScale)

		if material.transparencyEnabled {
			gl.Uniform1i(material.shader.GetUniform("transparencyEnabled"), 1)
			gl.ActiveTexture(gl.TEXTURE2)
			gl.BindTexture(gl.TEXTURE_2D, *material.transparencyTexture)
			gl.Uniform1i(material.shader.GetUniform("transparenctMap"), 2)
		} else {
			gl.Uniform1i(material.shader.GetUniform("transparencyEnabled"), 0)
		}

		gl.Uniform1f(material.shader.GetUniform("darkness"), darkness)
		gl.Uniform1f(material.shader.GetUniform("transparency"), material.transparency)

	case SHADER_TEXTURE:
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *material.texture)
		gl.Uniform3fv(material.shader.GetUniform("materialType"), 1, &SHADER_TEXTURE_UNI[0])
		gl.Uniform1i(material.shader.GetUniform("diffuseMap"), 0)
		gl.Uniform1i(material.shader.GetUniform("cubeDiffuseMap"), 1)
		gl.Uniform1f(material.shader.GetUniform("shine"), material.shine)
		gl.Uniform1f(material.shader.GetUniform("textureScale"), material.textureScale)

		if material.transparencyEnabled {
			gl.Uniform1i(material.shader.GetUniform("transparencyEnabled"), 1)
			gl.ActiveTexture(gl.TEXTURE2)
			gl.BindTexture(gl.TEXTURE_2D, *material.transparencyTexture)
			gl.Uniform1i(material.shader.GetUniform("transparencyMap"), 2)
		} else {
			gl.Uniform1i(material.shader.GetUniform("transparencyEnabled"), 0)
		}

		gl.Uniform1f(material.shader.GetUniform("darkness"), darkness)
		gl.Uniform1f(material.shader.GetUniform("transparency"), material.transparency)

		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)
		gl.EnableVertexAttribArray(2)

	case SHADER_CUBEMAP:
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, *material.texture)
		gl.Uniform3fv(material.shader.GetUniform("materialType"), 1, &SHADER_CUBEMAP_UNI[0])
		gl.Uniform1i(material.shader.GetUniform("cubeDiffuseMap"), 1)
		gl.Uniform1f(material.shader.GetUniform("shine"), material.shine)
		gl.Uniform1f(material.shader.GetUniform("darkness"), darkness)
		gl.Uniform1f(material.shader.GetUniform("transparency"), material.transparency)

		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)
		gl.EnableVertexAttribArray(2)
	}
}

func (material *Material) BecomeColor(rgb [3]float32) {
	material.shaderType = SHADER_COLOR
	material.color = [3]float32{rgb[0] / 255, rgb[1] / 255, rgb[2] / 255}
}

func (material *Material) BecomeTexture(t *uint32) {
	material.shaderType = SHADER_TEXTURE
	material.texture = t
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *t)
	gl.Uniform1i(material.shader.GetUniform("texture0"), int32(0))
}

func (material *Material) BecomeCubemap(c *uint32) {
	material.shaderType = SHADER_CUBEMAP
	material.texture = c
}

func (material *Material) AttachShader(s *ShaderProgram) {
	material.shader = s
}

func (material *Material) AttachTransparency(texture *uint32) {
	material.transparencyEnabled = true
	material.transparencyTexture = texture
}

func (material *Material) RemoveTransparency() {
	material.transparencyEnabled = false
	material.transparencyTexture = nil
}

func (materal *Material) GetColor() [3]float32 {
	return materal.color
}

func (material *Material) GetTexture() *uint32 {
	return material.texture
}

func (material *Material) GetShader() *ShaderProgram {
	return material.shader
}

func (material *Material) SetTextureScale(scale float32) {
	material.textureScale = scale
}

func (material *Material) EnableAnimation() {
	material.animationEnabled = true
	material.animationPlaying = ""
}

func (material *Material) AddFrame(frame *uint32, anim string) {
	material.animationTextures[anim] = append(material.animationTextures[anim], frame)
}

func (material *Material) SetAnimationFPS(anim string, s float64) {
	material.animationFPS[anim] = s
}

func (material *Material) PlayAnimation(anim string) {
	material.animationPlaying = anim
}

func (material *Material) SetTransparency(t float32) {
	material.transparency = t
}
