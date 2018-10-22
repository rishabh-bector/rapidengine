package cmd

import (
	"fmt"
	"net/http"
	"rapidengine/camera"
	"rapidengine/child"
	"rapidengine/configuration"
	"rapidengine/control"
	"rapidengine/input"
	"rapidengine/lighting"
	"rapidengine/terrain"
	"rapidengine/ui"

	"github.com/go-gl/mathgl/mgl32"
)

type Engine struct {
	Renderer   Renderer
	RenderFunc func(renderer *Renderer, inputs *input.Input)

	CollisionControl control.CollisionControl
	TextureControl   control.TextureControl
	InputControl     control.InputControl
	ShaderControl    control.ShaderControl
	LightControl     control.LightControl
	UIControl        control.UIControl

	TextControl control.TextControl
	FPSBox      *ui.TextBox

	FrameCount int

	Config configuration.EngineConfig
}

func NewEngine(config configuration.EngineConfig, renderFunc func(*Renderer, *input.Input)) Engine {
	e := Engine{
		Renderer:         NewRenderer(getEngineCamera(config.Dimensions, &config), &config),
		CollisionControl: control.NewCollisionControl(&config),
		TextureControl:   control.NewTextureControl(),
		InputControl:     control.NewInputControl(),
		ShaderControl:    control.NewShaderControl(),
		LightControl:     control.NewLightControl(),
		UIControl:        control.NewUIControl(),
		TextControl:      control.NewTextControl(&config),
		Config:           config,
		FrameCount:       0,
		RenderFunc:       renderFunc,
	}

	if e.Config.Profiling {
		http.HandleFunc("/", profileEndpoint)
		go http.ListenAndServe(":8080", nil)
	}

	if e.Config.ShowFPS {
		e.TextControl.LoadFont("../rapidengine/assets/fonts/roboto.ttf", "roboto", 64, 10)
		e.FPSBox = e.TextControl.NewTextBox("Rapid Engine", "roboto", float32(e.Config.ScreenWidth/2-100), float32(e.Config.ScreenHeight/2-50), 0.5, [3]float32{1, 1, 1})
		e.TextControl.AddTextBox(e.FPSBox)
	}

	e.ShaderControl.Initialize()
	e.Renderer.Initialize(&e)
	e.Renderer.AttachCallback(e.Update)

	if e.Config.Dimensions == 2 {
		l := lighting.NewDirectionLight(
			e.ShaderControl.GetShader("colorLighting"),
			[]float32{0, 0, 0},
			[]float32{0, 0, 0},
			[]float32{0, 0, 0},
			[]float32{0, 0, -1},
		)
		e.LightControl.SetDirectionalLight(&l)
	}
	if e.Config.Dimensions == 3 {
		l := lighting.NewDirectionLight(
			e.ShaderControl.GetShader("colorLighting"),
			[]float32{0.3, 0.3, 0.3},
			[]float32{0.8, 0.8, 0.8},
			[]float32{0, 0, 0},
			[]float32{1, -1, 1},
		)
		e.LightControl.SetDirectionalLight(&l)

		terrain.NewSkyBox("TropicalSunnyDay", &e.ShaderControl, &e.TextureControl, &e.Config)
	}

	return e
}

func NewEngineConfig(
	ScreenWidth,
	ScreenHeight,
	Dimensions int,
) configuration.EngineConfig {
	return configuration.NewEngineConfig(ScreenWidth, ScreenHeight, Dimensions)
}

func (engine *Engine) Initialize() {
	engine.Renderer.PreRenderChildren()
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
		engine.FPSBox.SetText(fmt.Sprintf("FPS: %v", int(1/renderer.DeltaFrameTime)))
		engine.FrameCount = 0
	}

	// Update controllers
	engine.LightControl.Update(x, y, z)
	engine.CollisionControl.Update(x, y, inputs)
	engine.UIControl.Update(inputs)
	engine.TextControl.Update()

	engine.FrameCount++
}

func (engine *Engine) NewChild2D() child.Child2D {
	c := NewChild2D(&engine.Config, &engine.CollisionControl)
	c.AttachShader(engine.Renderer.ShaderProgram)
	return c
}

func (engine *Engine) NewChild3D() child.Child3D {
	c := NewChild3D(&engine.Config, &engine.CollisionControl)
	c.AttachShader(engine.Renderer.ShaderProgram)
	return c
}

func (engine *Engine) StartRenderer() {
	if engine.Config.CollisionLines {
	}
	engine.Renderer.StartRenderer()
}

func (engine *Engine) Instance(c child.Child) {
	engine.Renderer.Instance(c)
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
