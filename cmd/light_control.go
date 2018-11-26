package cmd

import "rapidengine/lighting"

type LightControl struct {
	lightingEnabled    map[int]bool
	directionalEnabled map[int]bool

	DirLight      map[int]*lighting.DirectionLight
	pointLightMap map[int]*lighting.PointLight
}

func NewLightControl() LightControl {
	l := LightControl{
		lightingEnabled:    make(map[int]bool),
		directionalEnabled: make(map[int]bool),
		pointLightMap:      make(map[int]*lighting.PointLight),
		DirLight:           make(map[int]*lighting.DirectionLight),
	}
	l.EnableDirectionalLighting()
	return l
}

func (lightControl *LightControl) Update(cx, cy, cz float32) {
	if lightControl.lightingEnabled[0] {
		if lightControl.DirLight[0] != nil && lightControl.directionalEnabled[0] {
			lightControl.DirLight[0].UpdateShader(cx, cy, cz)
		}
		for ind, light := range lightControl.pointLightMap {
			light.UpdateShader(cx, cy, cz, ind)
		}
	}
}

func (lightControl *LightControl) PreRender() {
	if lightControl.lightingEnabled[0] {
		if lightControl.DirLight[0] != nil && lightControl.directionalEnabled[0] {
			lightControl.DirLight[0].PreRender()
		}
		for _, light := range lightControl.pointLightMap {
			light.PreRender()
		}
	}
}

func (lightControl *LightControl) InstanceLight(l *lighting.PointLight, ind int) {
	lightControl.pointLightMap[ind] = l
}

func (lightControl *LightControl) SetDirectionalLight(light *lighting.DirectionLight) {
	lightControl.DirLight[0] = light
}

func (lightControl *LightControl) EnableLighting() {
	lightControl.lightingEnabled[0] = true
}

func (lightControl *LightControl) DisableLighting() {
	lightControl.lightingEnabled[0] = false
}

func (lightControl *LightControl) EnableDirectionalLighting() {
	lightControl.directionalEnabled[0] = true
}

func (lightControl *LightControl) DisableDirectionalLighting() {
	lightControl.directionalEnabled[0] = false
}
