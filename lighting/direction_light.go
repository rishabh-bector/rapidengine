package lighting

import (
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DirectionLight struct {
	Ambient  []float32
	Diffuse  []float32
	Specular []float32

	Direction []float32
}

func NewDirectionLight(a, d, s, dir []float32) DirectionLight {
	return DirectionLight{
		Ambient:   a,
		Diffuse:   d,
		Specular:  s,
		Direction: dir,
	}
}

func (light *DirectionLight) PreRender() {

}

func (light *DirectionLight) UpdateShader(cx, cy, cz float32, shader *material.ShaderProgram) {
	c := []float32{cx, cy, cz}
	shader.Bind()
	gl.Uniform3fv(
		shader.GetUniform("dirLight.direction"),
		1, &light.Direction[0],
	)

	gl.Uniform3fv(
		shader.GetUniform("dirLight.ambient"),
		1, &light.Ambient[0],
	)

	gl.Uniform3fv(
		shader.GetUniform("dirLight.diffuse"),
		1, &light.Diffuse[0],
	)

	gl.Uniform3fv(
		shader.GetUniform("dirLight.specular"),
		1, &light.Specular[0],
	)

	gl.Uniform3fv(
		shader.GetUniform("viewPos"),
		1, &c[0],
	)
}
