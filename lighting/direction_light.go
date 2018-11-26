package lighting

import (
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DirectionLight struct {
	shader *material.ShaderProgram

	Ambient  []float32
	Diffuse  []float32
	Specular []float32

	Direction []float32
}

func NewDirectionLight(p *material.ShaderProgram, a, d, s, dir []float32) DirectionLight {
	return DirectionLight{
		shader:    p,
		Ambient:   a,
		Diffuse:   d,
		Specular:  s,
		Direction: dir,
	}
}

func (light *DirectionLight) PreRender() {
	light.shader.Bind()
}

func (light *DirectionLight) UpdateShader(cx, cy, cz float32) {
	c := []float32{cx, cy, cz}
	light.shader.Bind()
	gl.Uniform3fv(
		light.shader.GetUniform("dirLight.direction"),
		1, &light.Direction[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("dirLight.ambient"),
		1, &light.Ambient[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("dirLight.diffuse"),
		1, &light.Diffuse[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("dirLight.specular"),
		1, &light.Specular[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("viewPos"),
		1, &c[0],
	)
}
