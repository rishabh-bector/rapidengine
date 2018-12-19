package cmd

import (
	"fmt"
	"net/http"
	"rapidengine/camera"
	"rapidengine/configuration"
	"rapidengine/input"
	"rapidengine/lighting"
	"rapidengine/material"
	"rapidengine/ui"

	"github.com/sirupsen/logrus"

	"github.com/go-gl/mathgl/mgl32"
)

type Engine struct {
	Renderer   Renderer
	RenderFunc func(renderer *Renderer, inputs *input.Input)

	ChildControl     ChildControl
	SceneControl     SceneControl
	CollisionControl CollisionControl
	TextureControl   TextureControl
	MaterialControl  MaterialControl
	InputControl     InputControl
	ShaderControl    ShaderControl
	LightControl     LightControl
	UIControl        UIControl
	TerrainControl   TerrainControl
	TextControl      TextControl
	PostControl      PostControl

	FPSBox     *ui.TextBox
	FrameCount int

	Config *configuration.EngineConfig

	Logger *logrus.Logger
}

func NewEngine(config *configuration.EngineConfig, renderFunc func(*Renderer, *input.Input)) *Engine {
	e := Engine{
		// Main renderer
		Renderer: NewRenderer(getEngineCamera(config.Dimensions, config), config),

		// Package Controls
		ChildControl:     NewChildControl(),
		SceneControl:     NewSceneControl(),
		CollisionControl: NewCollisionControl(config),
		TextureControl:   NewTextureControl(config),
		InputControl:     NewInputControl(),
		ShaderControl:    NewShaderControl(),
		MaterialControl:  NewMaterialControl(),
		LightControl:     NewLightControl(),
		TerrainControl:   NewTerrainControl(),
		UIControl:        NewUIControl(),
		TextControl:      NewTextControl(config),
		PostControl:      NewPostControl(),

		// Configuration
		Config:     config,
		FrameCount: 0,

		// User render function
		RenderFunc: renderFunc,

		Logger: config.Logger,
	}

	e.ChildControl.Initialize(&e)
	e.SceneControl.Initialize(&e)
	e.ShaderControl.Initialize()
	e.MaterialControl.Initialize(&e)
	e.UIControl.Initialize(&e)
	e.TextControl.Initialize(&e)
	e.TerrainControl.Initialize(&e)
	e.CollisionControl.Initialize(&e)
	e.PostControl.Initialize(&e)
	e.LightControl.Initialize(&e)

	e.LightControl.Shaders = []*material.ShaderProgram{
		e.ShaderControl.GetShader("standard"),
		e.ShaderControl.GetShader("terrain"),
		e.ShaderControl.GetShader("foliage"),
		e.ShaderControl.GetShader("water"),
	}

	e.Renderer.Initialize(&e)
	e.Renderer.AttachCallback(e.Update)

	e.TextControl.LoadFont("../rapidengine/assets/fonts/avenir-next-regular.ttf", "avenir", 32, 0)
	if e.Config.ShowFPS {
		e.FPSBox = e.TextControl.NewTextBox("Rapid Engine", "roboto", float32(e.Config.ScreenWidth/2-100), float32(e.Config.ScreenHeight/2-50), 0.5, [3]float32{1, 1, 1})
		e.SceneControl.GetCurrentScene().InstanceText(e.FPSBox)
	}

	if e.Config.Dimensions == 3 {
		l := lighting.NewDirectionLight(
			[]float32{0.2, 0.2, 0.2},
			[]float32{0.9, 0.9, 0.9},
			[]float32{0.4, 0.4, 0.4},
			[]float32{1, 0, 0},
		)

		e.LightControl.SetDirectionalLight(&l)

		e.Renderer.SkyBoxEnabled = true
		e.Renderer.SkyBox = e.TerrainControl.NewSkyBox("TropicalSunnyDay", &e.ShaderControl, &e.TextureControl, e.Config)
	}

	return &e
}

func NewEngineConfig(
	ScreenWidth,
	ScreenHeight,
	Dimensions int,
) configuration.EngineConfig {
	return configuration.NewEngineConfig(ScreenWidth, ScreenHeight, Dimensions)
}

func (engine *Engine) Initialize() {
	engine.SceneControl.PreRenderChildren()
}

func (engine *Engine) Update(renderer *Renderer) {
	// Get camera position
	x, y, z := renderer.MainCamera.GetPosition()

	// Get user inputs
	inputs := engine.InputControl.Update(renderer.Window)

	// Call user frame function
	engine.RenderFunc(renderer, inputs)

	// Update FPS
	if engine.Config.ShowFPS && engine.FrameCount > 10 {
		engine.FPSBox.Text = fmt.Sprintf("FPS: %v", int(1/renderer.DeltaFrameTime))
		engine.FrameCount = 0
	}

	// Update controllers
	engine.TerrainControl.Update()
	engine.LightControl.Update(x, y, z)
	engine.CollisionControl.Update(x, y, inputs)
	engine.UIControl.Update(inputs)
	engine.TextControl.Update()

	engine.FrameCount++
}

func (engine *Engine) StartRenderer() {
	if engine.Config.CollisionLines {
	}
	engine.Renderer.StartRenderer()
}

func (engine *Engine) InstanceLight(l *lighting.PointLight) {
	engine.LightControl.InstanceLight(l, 0)
}

func (engine *Engine) SetDirectionalLight(l *lighting.DirectionLight) {
	engine.LightControl.SetDirectionalLight(l)
}

func (engine *Engine) EnableLighting() {
	engine.LightControl.EnableLighting()
}

func (engine *Engine) DisableLighting() {
	engine.LightControl.DisableLighting()
}

func (engine *Engine) Done() chan bool {
	return engine.Renderer.Done
}

func getEngineCamera(dimension int, config *configuration.EngineConfig) camera.Camera {
	if dimension == 2 {
		return camera.NewCamera2D(mgl32.Vec3{0, 0, 0}, float32(0.05), config)
	}
	if dimension == 3 {
		return camera.NewCamera3D(mgl32.Vec3{0, 0, 0}, float32(0.05), config)
	}
	return nil
}

func profileEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
