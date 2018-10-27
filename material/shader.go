package material

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

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

func (shaderProgram *ShaderProgram) GetID() uint32 {
	return shaderProgram.id
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
