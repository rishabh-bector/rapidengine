package rapidengine

import (
	"rapidengine/camera"
	"rapidengine/configuration"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Engine struct {
	Renderer   Renderer
	RenderFunc func(renderer *Renderer)

	CollisionControl CollisionControl
	TextureControl   TextureControl

	Shaders Shaders

	Config configuration.EngineConfig
}

func NewEngine(config configuration.EngineConfig, renderFunc func(*Renderer)) Engine {
	e := Engine{
		Renderer:         NewRenderer(camera.NewCamera2D(mgl32.Vec3{0, 0, 0}, float32(0.05), &config), &config),
		CollisionControl: NewCollisionControl(),
		TextureControl:   NewTextureControl(),
		Shaders:          NewShaders(),
		Config:           config,
		RenderFunc:       renderFunc,
	}
	e.Renderer.AttachCallback(e.Update)
	return e
}

func NewEngineConfig(
	ScreenWidth,
	ScreenHeight,
	Dimensions int,
) configuration.EngineConfig {
	return configuration.NewEngineConfig(ScreenWidth, ScreenHeight, Dimensions)
}

func (engine *Engine) Initialize() error {
	err := engine.Shaders.CompileShaders()
	if err != nil {
		engine.Config.Logger.Fatal(err)
		return err
	}
	gl.UseProgram(engine.Renderer.ShaderProgram)
	return nil
}

func (engine *Engine) InitializeRenderer() {
	engine.Renderer.PreRenderChildren()
}

func (engine *Engine) Update(renderer *Renderer) {
	x, y := renderer.MainCamera.GetPosition()
	engine.RenderFunc(renderer)
	engine.CollisionControl.Update(x, y)
}

func (engine *Engine) NewChild2D() Child2D {
	c := NewChild2D(&engine.Config, &engine.CollisionControl)
	c.AttachShader(engine.Renderer.ShaderProgram)
	return c
}

func (engine *Engine) StartRenderer() {
	if engine.Config.CollisionLines {
		// TODO: Add visible collision lines
	}
	engine.Renderer.StartRenderer()
}

func (engine *Engine) Instance(c Child) {
	engine.Renderer.Instance(c)
}

func (engine *Engine) Done() chan bool {
	return engine.Renderer.Done
}
