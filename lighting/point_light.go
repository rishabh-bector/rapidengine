package lighting

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type PointLight struct {
	program uint32

	ambient  []float32
	diffuse  []float32
	specular []float32

	constant  float32
	linear    float32
	quadratic float32

	position []float32
}

func NewPointLight(p uint32, a, d, s []float32, c, l, q float32) PointLight {
	return PointLight{
		program: p,

		ambient:  a,
		diffuse:  d,
		specular: s,

		constant:  c,
		linear:    l,
		quadratic: q,
	}
}

func (light *PointLight) PreRender() {}

func (light *PointLight) UpdateShader(cx, cy, cz float32, ind int) {
	c := []float32{cx, cy, cz}
	gl.UseProgram(light.program)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].ambient"+"\x00")),
		1, &light.ambient[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].diffuse"+"\x00")),
		1, &light.diffuse[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].specular"+"\x00")),
		1, &light.specular[0],
	)

	gl.Uniform1f(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].constant"+"\x00")),
		light.constant,
	)

	gl.Uniform1f(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].linear"+"\x00")),
		light.linear,
	)

	gl.Uniform1f(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].quadratic"+"\x00")),
		light.quadratic,
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("pointLights["+string(ind)+"].position"+"\x00")),
		1, &light.position[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("lmao.ambient"+"\x00")),
		1, &light.ambient[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("lmao.diffuse"+"\x00")),
		1, &light.diffuse[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("lmao.specular"+"\x00")),
		1, &light.specular[0],
	)

	gl.Uniform1f(
		gl.GetUniformLocation(light.program, gl.Str("lmao.constant"+"\x00")),
		light.constant,
	)

	gl.Uniform1f(
		gl.GetUniformLocation(light.program, gl.Str("lmao.linear"+"\x00")),
		light.linear,
	)

	gl.Uniform1f(
		gl.GetUniformLocation(light.program, gl.Str("lmao.quadratic"+"\x00")),
		light.quadratic,
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("lmao.position"+"\x00")),
		1, &light.position[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("viewPos"+"\x00")),
		1, &c[0],
	)
}

func (light *PointLight) SetPosition(pos []float32) {
	light.position = pos
}

func (light *PointLight) SetLevels(linear, quad float32) {
	light.linear = linear
	light.quadratic = quad
}
