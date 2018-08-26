package rapidengine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

type EngineConfig struct {
	ScreenWidth  int
	ScreenHeight int

	WindowTitle    string
	PolygonLines   bool
	CollisionLines bool

	Dimensions int
}

type Engine struct {
	Renderer   Renderer
	RenderFunc func(renderer *Renderer)

	CollisionControl CollisionControl

	Shaders Shaders

	Config EngineConfig

	Logger *logrus.Logger
}

func NewEngine(config EngineConfig, renderFunc func(*Renderer)) Engine {
	e := Engine{
		Renderer:         NewRenderer(NewCamera2D(mgl32.Vec3{0, 0, 0}, float32(0.02)), &config),
		CollisionControl: NewCollisionControl(),
		Shaders:          NewShaders(),
		Config:           config,
		Logger:           logrus.New(),
		RenderFunc:       renderFunc,
	}
	e.Renderer.AttachCallback(e.Update)
	return e
}

func (engine *Engine) Initialize() error {
	err := engine.Shaders.CompileShaders()
	if err != nil {
		engine.Logger.Fatal(err)
		return err
	}
	gl.UseProgram(engine.Renderer.ShaderProgram)
	return nil
}

func (engine *Engine) Update(renderer *Renderer) {
	engine.CollisionControl.Update()
	engine.RenderFunc(renderer)
}

func (engine *Engine) NewChild2D() Child2D {
	c := NewChild2D(&engine.Config)
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

func (engine *Engine) Done() chan bool {
	return engine.Renderer.Done
}
