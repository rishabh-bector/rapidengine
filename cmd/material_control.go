package cmd

import (
	"rapidengine/material"
	"rapidengine/state"
)

type MaterialControl struct {
	engine *Engine
}

func NewMaterialControl() MaterialControl {
	return MaterialControl{}
}

func (mc *MaterialControl) Initialize(engine *Engine) {
	mc.engine = engine

	e := uint32(1000)
	state.BoundTexture0 = e
	state.BoundTexture1 = e
	state.BoundTexture2 = e
	state.BoundTexture3 = e
	state.BoundTexture4 = e
	state.BoundTexture5 = e
	state.BoundTexture6 = e
}

func (mc *MaterialControl) NewBasicMaterial() *material.BasicMaterial {
	return material.NewBasicMaterial(mc.engine.ShaderControl.GetShader("basic"))
}

func (mc *MaterialControl) NewStandardMaterial() *material.StandardMaterial {
	return material.NewStandardMaterial(mc.engine.ShaderControl.GetShader("standard"))
}

func (mc *MaterialControl) NewPBRMaterial() *material.PBRMaterial {
	return material.NewPBRMaterial(mc.engine.ShaderControl.GetShader("pbr"))
}

func (mc *MaterialControl) NewCubemapMaterial() *material.CubemapMaterial {
	return material.NewCubemapMaterial(mc.engine.ShaderControl.GetShader("skybox"))
}

func (mc *MaterialControl) NewTerrainMaterial() *material.TerrainMaterial {
	return material.NewTerrainMaterial(mc.engine.ShaderControl.GetShader("terrain"))
}

func (mc *MaterialControl) NewFoliageMaterial() *material.FoliageMaterial {
	return material.NewFoliageMaterial(mc.engine.ShaderControl.GetShader("foliage"))
}

func (mc *MaterialControl) NewWaterMaterial() *material.WaterMaterial {
	return material.NewWaterMaterial(mc.engine.ShaderControl.GetShader("water"))
}
