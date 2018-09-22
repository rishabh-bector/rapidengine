package rapidengine

import (
	"rapidengine/camera"
	"rapidengine/configuration"
	"rapidengine/input"

	"github.com/go-gl/mathgl/mgl32"
)

type Engine struct {
	Renderer   Renderer
	RenderFunc func(renderer *Renderer, inputs *input.Input)

	CollisionControl CollisionControl
	TextureControl   TextureControl
	InputControl     input.InputControl
	ShaderControl    ShaderControl
	LightControl     LightControl

	Config configuration.EngineConfig
}

func NewEngine(config configuration.EngineConfig, renderFunc func(*Renderer, *input.Input)) Engine {
	e := Engine{
		Renderer:         NewRenderer(getEngineCamera(config.Dimensions, &config), &config),
		CollisionControl: NewCollisionControl(),
		TextureControl:   NewTextureControl(),
		InputControl:     input.NewInputControl(),
		ShaderControl:    NewShaderControl(),
		LightControl:     NewLightControl(),
		Config:           config,
		RenderFunc:       renderFunc,
	}

	e.ShaderControl.Initialize()
	e.Renderer.AttachCallback(e.Update)

	if e.Config.Dimensions == 2 {
		l := NewDirectionLight(
			e.ShaderControl.GetShader("colorLighting"),
			[]float32{0, 0, 0},
			[]float32{0, 0, 0},
			[]float32{0, 0, 0},
			[]float32{0, 0, -1},
		)
		e.LightControl.SetDirectionalLight(&l)
	}
	if e.Config.Dimensions == 3 {
		l := NewDirectionLight(
			e.ShaderControl.GetShader("colorLighting"),
			[]float32{0.1, 0.1, 0.1},
			[]float32{0.6, 0.6, 0.6},
			[]float32{0.2, 0.2, 0.2},
			[]float32{1, -1, 1},
		)
		e.LightControl.SetDirectionalLight(&l)

		NewSkyBox(
			"../rapidengine/skybox/right.jpg",
			"../rapidengine/skybox/left.jpg",
			"../rapidengine/skybox/top.jpg",
			"../rapidengine/skybox/bottom.jpg",
			"../rapidengine/skybox/front.jpg",
			"../rapidengine/skybox/back.jpg", &e,
		)
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
	x, y, z := renderer.MainCamera.GetPosition()
	engine.RenderFunc(renderer, engine.InputControl.Update(renderer.Window))
	engine.LightControl.Update(x, y, z)
	engine.CollisionControl.Update(x, y)
}

func (engine *Engine) NewChild2D() Child2D {
	c := NewChild2D(&engine.Config, &engine.CollisionControl)
	c.AttachShader(engine.Renderer.ShaderProgram)
	return c
}

func (engine *Engine) NewChild3D() Child3D {
	c := NewChild3D(&engine.Config, &engine.CollisionControl)
	c.AttachShader(engine.Renderer.ShaderProgram)
	return c
}

func (engine *Engine) StartRenderer() {
	if engine.Config.CollisionLines {
	}
	engine.Renderer.StartRenderer()
}

func (engine *Engine) Instance(c Child) {
	engine.Renderer.Instance(c)
}

func (engine *Engine) InstanceLight(l *PointLight) {
	engine.LightControl.InstanceLight(l, 0)
}

func (engine *Engine) SetDirectionalLight(l *DirectionLight) {
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
