package cmd

import "rapidengine/material"

type ShaderControl struct {
	programs map[string]*material.ShaderProgram
}

func NewShaderControl() ShaderControl {
	return ShaderControl{make(map[string]*material.ShaderProgram)}
}

func (shaderControl *ShaderControl) BindShader(name string) {
	shaderControl.programs[name].Bind()
}

func (shaderControl *ShaderControl) Initialize() {
	shaderControl.programs = map[string]*material.ShaderProgram{
		"basic":    &material.BasicProgram,
		"standard": &material.StandardProgram,
		"skybox":   &material.SkyBoxProgram,
	}
	for _, prog := range shaderControl.programs {
		prog.Compile()
	}
}

func (shaderControl *ShaderControl) GetShader(name string) *material.ShaderProgram {
	return shaderControl.programs[name]
}