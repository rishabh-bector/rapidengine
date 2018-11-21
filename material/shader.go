package material

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderProgram struct {
	id             uint32
	vertexShader   string
	fragmentShader string

	uniformLocations   map[string]int32
	attributeLocations map[string]uint32
}

func (shaderProgram *ShaderProgram) Bind() {
	b := shaderProgram.id
	gl.UseProgram(b)
}

func (shaderProgram *ShaderProgram) RebindAttribLocations() {
	for attrib, location := range shaderProgram.attributeLocations {
		gl.BindAttribLocation(shaderProgram.id, location, gl.Str(attrib+"\x00"))
	}
}

func (shaderProgram *ShaderProgram) GetUniform(name string) int32 {
	return shaderProgram.uniformLocations[name]
}

func (shaderProgram *ShaderProgram) GetID() uint32 {
	return shaderProgram.id
}

func (shaderProgram *ShaderProgram) Compile() {
	vert, err := ioutil.ReadFile(shaderProgram.vertexShader)
	if err != nil {
		panic(err)
	}
	vertexShader, err := CompileShader(string(vert)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	frag, err := ioutil.ReadFile(shaderProgram.fragmentShader)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(string(frag)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram.id = gl.CreateProgram()
	gl.AttachShader(shaderProgram.id, vertexShader)
	gl.AttachShader(shaderProgram.id, fragmentShader)
	gl.LinkProgram(shaderProgram.id)

	for uni := range shaderProgram.uniformLocations {
		shaderProgram.uniformLocations[uni] = gl.GetUniformLocation(shaderProgram.id, gl.Str(uni+"\x00"))
	}

	for attrib, location := range shaderProgram.attributeLocations {
		gl.BindAttribLocation(shaderProgram.id, location, gl.Str(attrib+"\x00"))
	}
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

//  --------------------------------------------------
//  Shader Programs
//  --------------------------------------------------

var BasicProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/glsl/basic/basic.vert",
	fragmentShader: "../rapidengine/material/glsl/basic/basic.frag",
	uniformLocations: map[string]int32{
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,

		"diffuseLevel":    0,
		"hue":             0,
		"darkness":        0,
		"diffuseMap":      0,
		"diffuseMapScale": 0,
		"alphaMapLevel":   0,
		"alphaMap":        0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      1,
	},
}

var StandardProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/glsl/standard/standard.vert",
	fragmentShader: "../rapidengine/material/glsl/standard/standard.frag",
	uniformLocations: map[string]int32{
		"textureScale":        0,
		"copyingEnabled":      0,
		"transparency":        0,
		"modelMtx":            0,
		"viewMtx":             0,
		"projectionMtx":       0,
		"materialType":        0,
		"diffuseMap":          0,
		"cubeDiffuseMap":      0,
		"darkness":            0,
		"color":               0,
		"shine":               0,
		"transparencyEnabled": 0,
		"transparencyMap":     0,
		"viewPos":             0,

		"dirLight.direction": 0,
		"dirLight.ambient":   0,
		"dirLight.diffuse":   0,
		"dirLight.specular":  0,

		"pointLights": 0,
		"lmao":        0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
		"tex":      1,
		"normal":   2,
	},
}

var SkyBoxProgram = ShaderProgram{
	vertexShader:   "../rapidengine/material/glsl/skybox/skybox.vert",
	fragmentShader: "../rapidengine/material/glsl/skybox/skybox.frag",
	uniformLocations: map[string]int32{
		"color":         0,
		"transparency":  0,
		"modelMtx":      0,
		"viewMtx":       0,
		"projectionMtx": 0,
		"darkness":      0,
	},
	attributeLocations: map[string]uint32{
		"position": 0,
	},
}
