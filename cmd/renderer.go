package cmd

//   --------------------------------------------------
//   Render.go contains the main render loop, as well as
//   functions to initialize OpenGL and GLFW. A renderer
//   has a list of "children" which it renders every frame.
//   --------------------------------------------------

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/pkg/profile"
	log "github.com/sirupsen/logrus"

	"rapidengine/camera"
	"rapidengine/child"
	"rapidengine/configuration"
	"rapidengine/input"
	"rapidengine/material"
	"rapidengine/terrain"
)

// Renderer contains the information required for
// the main render loop
type Renderer struct {
	// GLFW Window object
	Window *glfw.Window

	// Current shader program
	ShaderProgram uint32

	// Currently bound VAO
	CurrentBoundVAO    uint32
	CurrentBoundShader uint32

	// Per-frame callback from the user
	RenderFunc func(renderer *Renderer)

	// Scene Camera
	MainCamera camera.Camera

	// Current camera position
	camX float32
	camY float32
	camZ float32

	// Render Distance
	RenderDistance float32

	// Skybox
	SkyBoxEnabled bool
	SkyBox        *terrain.SkyBox

	// Engine Configuration
	Config *configuration.EngineConfig

	// Default Material
	DefaultMaterial1 *material.BasicMaterial
	DefaultMaterial2 *material.BasicMaterial

	// FrameTime
	DeltaFrameTime float64
	LastFrameTime  float64
	MinFrameTime   float64
	TotalFrameTime float64

	// Termination Channel
	Done chan bool

	// Engine
	engine *Engine
}

// StartRenderer contains the main render loop
func (renderer *Renderer) StartRenderer() {
	if renderer.Config.Profiling {
		defer profile.Start().Stop()
	}

	gl.ClearColor(float32(0)/255, float32(0)/255, float32(0)/255, 1)

	// Render loop
	for !renderer.Window.ShouldClose() {
		renderer.renderFrame()
	}

	renderer.Config.Logger.Info("Terminating...")
	glfw.Terminate()
	renderer.Done <- true
}

// RenderFrame renders a single frame to the screen
func (renderer *Renderer) renderFrame() {
	renderer.engine.PostControl.UpdateFrameBuffers()

	// Render skybox
	if renderer.SkyBoxEnabled {
		renderer.SkyBox.Render(renderer.MainCamera)
	}

	// Render children
	renderer.RenderChildren()

	// Call user render loop
	renderer.RenderFunc(renderer)

	// Update camera
	renderer.MainCamera.Look(renderer.DeltaFrameTime)
	renderer.camX, renderer.camY, renderer.camZ = renderer.MainCamera.GetPosition()

	// Post processing update
	renderer.engine.PostControl.Update()

	// Update window buffers
	renderer.Window.SwapBuffers()

	// Frame logic
	renderer.TotalFrameTime = glfw.GetTime()
	renderer.DeltaFrameTime = renderer.TotalFrameTime - renderer.LastFrameTime
	renderer.LastFrameTime = renderer.TotalFrameTime

	if renderer.DeltaFrameTime < renderer.MinFrameTime {
		time.Sleep(time.Duration(1000000000 * (renderer.MinFrameTime - renderer.DeltaFrameTime)))
		renderer.DeltaFrameTime = renderer.MinFrameTime
	}
}

// ForceUpdate forces a frame render
func (renderer *Renderer) ForceUpdate() {
	renderer.engine.PostControl.UpdateFrameBuffers()
	renderer.RenderChildren()
	renderer.engine.PostControl.Update()
	renderer.engine.TextControl.Update()

	renderer.Window.SwapBuffers()
}

// RenderChildren binds the appropriate shaders and Vertex Array for each child,
// or child copy, and draws them to the screen using an element buffer
func (renderer *Renderer) RenderChildren() {
	if renderer.engine.SceneControl.GetCurrentScene().IsAutomaticRendering() {
		for _, child := range renderer.engine.SceneControl.GetCurrentChildren() {
			go child.RemoveCurrentCopies()
			if !child.CheckCopyingEnabled() {
				renderer.RenderChild(child)
			} else {
				renderer.RenderChildCopies(child)
			}
		}
	}
}

// RenderChild renders a single child to the screen
func (renderer *Renderer) RenderChild(c child.Child) {
	c.Update(renderer.MainCamera, renderer.DeltaFrameTime, renderer.TotalFrameTime)
}

// RenderChildCopies renders all copies of a child
func (renderer *Renderer) RenderChildCopies(c child.Child) {
	renderer.BindChild(c)

	copies := *(c.GetCopies())
	for x := 0; x < c.GetNumCopies(); x++ {
		renderer.RenderCopy(c, copies[x])
	}
}

// RenderCopy renders a single copy of a child
func (renderer *Renderer) RenderCopy(c child.Child, cpy child.ChildCopy) {
	renderer.BindChild(c)

	if renderer.Config.Dimensions == 2 {
		if (c.GetSpecificRenderDistance() != 0 && InBounds2D(cpy.X, cpy.Y, float32(renderer.camX), float32(renderer.camY), c.GetSpecificRenderDistance())) ||
			InBounds2D(cpy.X, cpy.Y, float32(renderer.camX), float32(renderer.camY), renderer.RenderDistance) {

			c.RenderCopy(cpy, renderer.MainCamera)
			c.AddCurrentCopy(cpy)
		}
	}

	if renderer.Config.Dimensions == 3 {
		if InBounds3D(cpy.X, cpy.Y, cpy.Z, float32(renderer.camX), float32(renderer.camY), float32(renderer.camZ), renderer.RenderDistance) {
			c.RenderCopy(cpy, renderer.MainCamera)

			c.AddCurrentCopy(cpy)
		}
	}
}

// BindChild intelligently binds the VAO & Shader of a child
func (renderer *Renderer) BindChild(c child.Child) {
	c.BindChild()
}

// InBounds2D checks if a particular x/y is within the given render distance
func InBounds2D(x, y, camX, camY, renderDistance float32) bool {
	if x < camX+renderDistance &&
		x > camX-renderDistance &&
		y < camY+renderDistance &&
		y > camY-renderDistance {
		return true
	}
	return false
}

// InBounds3D checks if a particular x/y/z is within the given render distance
func InBounds3D(x, y, z, camX, camY, camZ, renderDistance float32) bool {
	if x < camX+renderDistance &&
		x > camX-renderDistance &&
		y < camY+renderDistance &&
		y > camY-renderDistance &&
		z < camZ+renderDistance &&
		z > camZ-renderDistance {
		return true
	}
	return false
}

// NewRenderer creates a new renderer, and takes in a renderFunc which
// is called every frame, allowing the User to have frame-by-frame control
func NewRenderer(camera camera.Camera, config *configuration.EngineConfig) Renderer {
	win := InitGLFW(config)
	s := uint32(0)
	r := Renderer{
		Window:         win,
		ShaderProgram:  s,
		RenderFunc:     func(r *Renderer) {},
		RenderDistance: 1000,
		MinFrameTime:   1 / float64(config.MaxFPS),
		Done:           make(chan bool),
		MainCamera:     camera,
		Config:         config,
	}

	r.Window.SetCursorPosCallback(input.MouseCallback)
	r.Window.SetMouseButtonCallback(input.MouseButtonCallback)
	r.Window.SetScrollCallback(input.ScrollCallback)

	gl.Init()

	return r
}

func (renderer *Renderer) Initialize(engine *Engine) {
	renderer.engine = engine

	engine.TextureControl.NewTexture("../rapidengine/assets/abstract.jpg", "default", "linear")

	dm1 := renderer.engine.MaterialControl.NewBasicMaterial()
	dm1.Hue = [4]float32{46, 49, 49, 255}

	dm2 := renderer.engine.MaterialControl.NewBasicMaterial()
	dm2.Hue = [4]float32{211, 84, 0, 255}

	renderer.DefaultMaterial1 = dm1
	renderer.DefaultMaterial2 = dm2
}

// AttachCallback attaches a callback function to the renderer,
// to be called per-frame
func (renderer *Renderer) AttachCallback(f func(*Renderer)) {
	renderer.RenderFunc = f
}

func InitGLFW(config *configuration.EngineConfig) *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	var m *glfw.Monitor
	if config.FullScreen {
		m = glfw.GetPrimaryMonitor()
	} else {
		m = nil
	}

	window, err := glfw.CreateWindow(config.ScreenWidth, config.ScreenHeight, config.WindowTitle, m, nil)
	if err != nil {
		log.Fatal(err)
	}

	if config.Dimensions == 3 {
		//window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}

	window.MakeContextCurrent()

	if config.AntiAliasing {
		glfw.WindowHint(glfw.Samples, 8)
	}

	if !config.VSync {
		glfw.SwapInterval(0)
	}

	return window
}

func InitOpenGL(config *configuration.EngineConfig) uint32 {
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("Using OpenGL Version ", version)

	if config.PolygonLines {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if config.Dimensions == 3 {
		gl.Enable(gl.DEPTH_TEST)
		gl.Disable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}

	if config.GammaCorrection {
		gl.Enable(gl.FRAMEBUFFER_SRGB)
	}

	if config.AntiAliasing {
		gl.Enable(gl.MULTISAMPLE)
	}

	return 0
}

func (renderer *Renderer) ResetOpenGL(config *configuration.EngineConfig) {
	if config.PolygonLines {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if config.Dimensions == 3 {
		gl.Enable(gl.DEPTH_TEST)
		gl.Disable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}

	if config.GammaCorrection {
		gl.Enable(gl.FRAMEBUFFER_SRGB)
	}

	if config.AntiAliasing {
		gl.Enable(gl.MULTISAMPLE)
	}
}

func (renderer *Renderer) EnablePolygonLines() {
	renderer.engine.Config.PolygonLines = true
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func (renderer *Renderer) DisablePolygonLines() {
	renderer.engine.Config.PolygonLines = false
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

// SetRenderDistance sets the render distance
func (renderer *Renderer) SetRenderDistance(distance float32) {
	renderer.RenderDistance = distance
}

func (renderer *Renderer) DisableCursor() {
	renderer.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (renderer *Renderer) EnableCursor() {
	renderer.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

// CheckError decodes the various unhelpful error codes
// which OpenGL sometimes creates
func CheckError(tag string) {
	if err := gl.GetError(); err != 0 {
		var errString = ""
		switch err {
		case 0:
			return
		case 1280:
			errString = "An Enumeration parameter is not legal"
		case 1281:
			errString = "A value parameter is not legal"
		case 1282:
			errString = "A state for a command is not legal for its given parameters"
		case 1283:
			errString = "A stack pushing operation caused a stack overflow"
		case 1284:
			errString = "A stack popping operation occurred when the stack was at its lowest point"
		case 1285:
			errString = "A memory allocation could not allocate enough memory"
		case 1286:
			errString = "Attempting to read/write from an incomplete framebuffer"
		default:
			errString = "Unknown error"
		}
		log.Error(tag, ": ", errString)
	}
}
