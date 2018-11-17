package lighting

import (
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DirectionLight struct {
	shader *material.ShaderProgram

	ambient  []float32
	diffuse  []float32
	specular []float32

	direction []float32
}

func NewDirectionLight(p *material.ShaderProgram, a, d, s, dir []float32) DirectionLight {
	return DirectionLight{
		shader:    p,
		ambient:   a,
		diffuse:   d,
		specular:  s,
		direction: dir,
	}
}

func (light *DirectionLight) PreRender() {
	light.shader.Bind()
}

func (light *DirectionLight) UpdateShader(cx, cy, cz float32) {
	c := []float32{cx, cy, cz}
	light.shader.Bind()
	gl.Uniform3fv(
		light.shader.GetUniform("dirlight.direction"),
		1, &light.direction[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("dirlight.ambient"),
		1, &light.ambient[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("dirlight.diffuse"),
		1, &light.diffuse[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("dirlight.specular"),
		1, &light.specular[0],
	)

	gl.Uniform3fv(
		light.shader.GetUniform("viewPos"),
		1, &c[0],
	)
}
