package lighting

import (
	"fmt"
	"rapidengine/material"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type PointLight struct {
	Ambient  []float32
	Diffuse  []float32
	specular []float32

	constant  float32
	linear    float32
	quadratic float32

	Position []float32
}

func NewPointLight(a, d, s []float32, c, l, q float32) *PointLight {
	return &PointLight{
		Ambient:  a,
		Diffuse:  d,
		specular: s,

		constant:  c,
		linear:    l,
		quadratic: q,
	}
}

func (light *PointLight) PreRender() {}

func (light *PointLight) UpdateShader(cx, cy, cz float32, ind int, shader *material.ShaderProgram) {
	c := []float32{cx, cy, cz}
	shader.Bind()

	gl.Uniform3fv(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].ambient"+"\x00")),
		1, &light.Ambient[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].diffuse"+"\x00")),
		1, &light.Diffuse[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].specular"+"\x00")),
		1, &light.specular[0],
	)

	gl.Uniform1f(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].constant"+"\x00")),
		light.constant,
	)

	gl.Uniform1f(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].linear"+"\x00")),
		light.linear,
	)

	gl.Uniform1f(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].quadratic"+"\x00")),
		light.quadratic,
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(shader.GetID(), gl.Str("pointLights["+fmt.Sprint(ind)+"].position"+"\x00")),
		1, &light.Position[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(shader.GetID(), gl.Str("viewPos"+"\x00")),
		1, &c[0],
	)
}

func (light *PointLight) SetPosition(pos []float32) {
	light.Position = pos
}

func (light *PointLight) SetLevels(linear, quad float32) {
	light.linear = linear
	light.quadratic = quad
}
