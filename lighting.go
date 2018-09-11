package rapidengine

import "github.com/go-gl/gl/v4.1-core/gl"

type LightControl struct {
	LightMap map[int]Light
}

func NewLightControl() LightControl {
	return LightControl{}
}

func (lightControl *LightControl) Update(cx, cy, cz float32) {
	for _, light := range lightControl.LightMap {
		light.UpdateShader(cx, cy, cz)
	}
}

func (lightControl *LightControl) Initialize() {
	lightControl.LightMap = make(map[int]Light)
}

func (lightControl *LightControl) PreRender() {
	for _, light := range lightControl.LightMap {
		light.PreRender()
	}
}

func (lightControl *LightControl) InstanceLight(l Light, ind int) {
	lightControl.LightMap[ind] = l
}

type Light struct {
	program uint32

	ambient  float32
	diffuse  float32
	specular float32

	x float32
	y float32
	z float32
}

func NewLight(p uint32, a, d, s, x, y, z float32) Light {
	return Light{p, a, d, s, x, y, z}
}

func (light *Light) PreRender() {
	gl.UseProgram(light.program)
}

func (light *Light) UpdateShader(cx, cy, cz float32) {
	l := []float32{light.ambient, light.diffuse, light.specular}
	p := []float32{light.x, light.y, light.z}
	c := []float32{cx, cy, cz}
	gl.UseProgram(light.program)
	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("light"+"\x00")),
		1, &l[0],
	)
	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("lightPos"+"\x00")),
		1, &p[0],
	)
	gl.Uniform3fv(
		gl.GetUniformLocation(light.program, gl.Str("viewPos"+"\x00")),
		1, &c[0],
	)
}
