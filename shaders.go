package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shaders struct {
	shaderList []string
	typeList   []uint32
	idList     []uint32
}

const VertexShaderSource = `

		#version 410

		in vec3 position;
		in vec2 tex;

		out vec3 vertexColor;
		out vec2 texCoord;
		
		void main() {
			gl_Position = vec4(position, 1.0);
			texCoord = tex;
		}
	
	` + "\x00"

const FragmentShaderSource = `

		#version 410

		in vec2 texCoord;

		out vec4 outColor;

		uniform sampler2D texture0;

		void main() {
			outColor = texture(texture0, texCoord);
		}
		
	` + "\x00"

func NewShaders() *Shaders {
	return &Shaders{
		shaderList: []string{
			VertexShaderSource,
			FragmentShaderSource,
		},
		typeList: []uint32{
			gl.VERTEX_SHADER,
			gl.FRAGMENT_SHADER,
		},
		idList: []uint32{},
	}
}

func (shaders *Shaders) CompileShaders() error {
	for i := 0; i < len(shaders.shaderList); i++ {
		s, err := CompileShader(shaders.shaderList[i], shaders.typeList[i])
		if err != nil {
			return err
		}
		shaders.idList = append(shaders.idList, s)
	}
	return nil
}

func (shaders *Shaders) CleanUp() {
	for _, ind := range shaders.idList {
		gl.DeleteShader(ind)
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
