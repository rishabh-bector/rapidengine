package cmd

import (
	"rapidengine/material"
)

type MaterialControl struct {
	engine *Engine
}

func NewMaterialControl() MaterialControl {
	return MaterialControl{}
}

func (mc *MaterialControl) Initialize(engine *Engine) {
	mc.engine = engine
}

func (mc *MaterialControl) NewBasicMaterial() *material.BasicMaterial {
	return material.NewBasicMaterial(mc.engine.ShaderControl.GetShader("basic"))
}

func (mc *MaterialControl) NewStandardMaterial() *material.StandardMaterial {
	return material.NewStandardMaterial(mc.engine.ShaderControl.GetShader("standard"))
}

func (mc *MaterialControl) NewCubemapMaterial() *material.CubemapMaterial {
	return material.NewCubemapMaterial(mc.engine.ShaderControl.GetShader("skybox"))
}
