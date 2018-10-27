package lighting

import "github.com/go-gl/gl/v4.1-core/gl"

type DirectionLight struct {
	program uint32

	ambient  []float32
	diffuse  []float32
	specular []float32

	direction []float32
}

func NewDirectionLight(p uint32, a, d, s, dir []float32) DirectionLight {
	return DirectionLight{
		program:   p,
		ambient:   a,
		diffuse:   d,
		specular:  s,
		direction: dir,
	}
}

func (light *DirectionLight) PreRender() {
	gl.UseProgram(light.program)
}

func (light *DirectionLight) UpdateShader(cx, cy, cz float32) {
	c := []float32{cx, cy, cz}
	gl.UseProgram(light.program)
	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("dirLight.direction"+"\x00")),
		1, &light.direction[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("dirLight.ambient"+"\x00")),
		1, &light.ambient[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("dirLight.diffuse"+"\x00")),
		1, &light.diffuse[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("dirLight.specular"+"\x00")),
		1, &light.specular[0],
	)

	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("viewPos"+"\x00")),
		1, &c[0],
	)
}
