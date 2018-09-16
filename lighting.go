package rapidengine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type LightControl struct {
	lightingEnabled map[int]bool
	directionLight  map[int]*DirectionLight
	pointLightMap   map[int]PointLight
}

func NewLightControl() LightControl {
	return LightControl{
		lightingEnabled: make(map[int]bool),
		pointLightMap:   make(map[int]PointLight),
		directionLight:  make(map[int]*DirectionLight),
	}
}

func (lightControl *LightControl) Update(cx, cy, cz float32) {
	if lightControl.lightingEnabled[0] {
		if lightControl.directionLight[0] != nil {
			lightControl.directionLight[0].UpdateShader(cx, cy, cz)
		}
		for ind, light := range lightControl.pointLightMap {
			light.UpdateShader(cx, cy, cz, ind)
		}
	}
}

func (lightControl *LightControl) PreRender() {
	if lightControl.lightingEnabled[0] {
		if lightControl.directionLight[0] != nil {
			lightControl.directionLight[0].PreRender()
		}
		for _, light := range lightControl.pointLightMap {
			light.PreRender()
		}
	}
}

func (lightControl *LightControl) InstanceLight(l PointLight, ind int) {
	lightControl.pointLightMap[ind] = l
}

func (lightControl *LightControl) SetDirectionLight(light *DirectionLight) {
	lightControl.directionLight[0] = light
}

func (lightControl *LightControl) EnableLighting() {
	lightControl.lightingEnabled[0] = true
}

func (lightControl *LightControl) DisableLighting() {
	lightControl.lightingEnabled[0] = false
}

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
		gl.GetUniformLocation(light.program, gl.Str("viewPos"+"\x00")),
		1, &c[0],
	)
}

func (light *PointLight) SetPosition(pos []float32) {
	light.position = pos
}
