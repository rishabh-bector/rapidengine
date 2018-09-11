package rapidengine

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderControl struct {
	programs map[string]*ShaderProgram
}

type ShaderProgram struct {
	id             uint32
	vertexShader   string
	fragmentShader string
}

func (shaderProgram *ShaderProgram) Compile() {
	vertexShader, err := CompileShader(shaderProgram.vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := CompileShader(shaderProgram.fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	shaderProgram.id = gl.CreateProgram()
	gl.AttachShader(shaderProgram.id, vertexShader)
	gl.AttachShader(shaderProgram.id, fragmentShader)
	gl.LinkProgram(shaderProgram.id)
}

func NewShaderControl() ShaderControl {
	return ShaderControl{make(map[string]*ShaderProgram)}
}

func (shaderControl *ShaderControl) Initialize() {
	shaderControl.programs = map[string]*ShaderProgram{
		"texture":       &TextureProgram,
		"colorLighting": &ColorLightingProgram,
		"color":         &ColorProgram,
	}
	for _, prog := range shaderControl.programs {
		prog.Compile()
	}
}

func (shaderControl *ShaderControl) GetShader(name string) uint32 {
	return shaderControl.programs[name].id
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
