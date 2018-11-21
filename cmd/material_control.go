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
