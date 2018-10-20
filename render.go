package rapidengine

//   --------------------------------------------------
//   Render.go contains the main render loop, as well as
//   functions to initialize OpenGL and GLFW. A renderer
//   has a list of "children" which it renders every frame.
//   --------------------------------------------------

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	log "github.com/sirupsen/logrus"

	"rapidengine/camera"
	"rapidengine/configuration"
	"rapidengine/input"
)

// Renderer contains the information required for
// the main render loop
type Renderer struct {
	// GLFW Window object
	Window *glfw.Window

	// Current shader program
	ShaderProgram uint32

	// Children to be rendered
	Children []Child

	// Render all children every frame (default)
	AutomaticRendering bool

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
	SkyBox        *SkyBox

	// Engine Configuration
	Config *configuration.EngineConfig

	// Default Material
	DefaultMaterial Material

	// FrameTime
	DeltaFrameTime float64
	LastFrameTime  float64

	// Termination Channel
	Done chan bool
}

// StartRenderer contains the main render loop
func (renderer *Renderer) StartRenderer() {
	if renderer.Config.Profiling {
		//defer profile.Start().Stop()
	}
	gl.ClearColor(float32(0)/255, float32(0)/255, float32(0)/255, 0.9)
	for !renderer.Window.ShouldClose() {

		// Clear screen buffers
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Render skybox
		if renderer.SkyBoxEnabled {
			renderer.SkyBox.Render(renderer.MainCamera)
		}

		// Render children
		if renderer.AutomaticRendering {
			renderer.RenderChildren()
		}

		// Call user render loop
		renderer.RenderFunc(renderer)

		// Update camera
		renderer.MainCamera.Look()
		renderer.camX, renderer.camY, renderer.camZ = renderer.MainCamera.GetPosition()

		// Update window buffers
		renderer.Window.SwapBuffers()

		// Frame logic
		currentFrame := glfw.GetTime()
		renderer.DeltaFrameTime = currentFrame - renderer.LastFrameTime
		renderer.LastFrameTime = currentFrame
		//println(int(renderer.DeltaFrameTime * 1000))
	}

	renderer.Config.Logger.Info("Terminating...")
	glfw.Terminate()
	renderer.Done <- true
}

// PreRenderChildren calls the PreRender method of each child,
// for initialization
func (renderer *Renderer) PreRenderChildren() {
	if renderer.Config.SingleMaterial {
		renderer.DefaultMaterial.PreRender()
	}
	for _, child := range renderer.Children {
		child.PreRender(renderer.MainCamera)
	}
}

// RenderChildren binds the appropriate shaders and Vertex Array for each child,
// or child copy, and draws them to the screen using an element buffer
func (renderer *Renderer) RenderChildren() {
	if renderer.Config.SingleMaterial {
		renderer.DefaultMaterial.Render(0, 1)
	}
	for _, child := range renderer.Children {
		go child.RemoveCurrentCopies()
		if !child.CheckCopyingEnabled() {
			renderer.RenderChild(child)
		} else {
			renderer.RenderChildCopies(child)
		}
	}
}

// RenderChild renders a single child to the screen
func (renderer *Renderer) RenderChild(child Child) {
	BindChild(child)

	child.Update(renderer.MainCamera, renderer.DeltaFrameTime, renderer.LastFrameTime)
	renderer.DrawChild(child)

	gl.BindVertexArray(0)
}

// DrawChild draws the child's vertices to the screen
func (renderer *Renderer) DrawChild(child Child) {
	gl.DrawElements(gl.TRIANGLES, child.GetNumVertices(), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// RenderChildCopies renders all copies of a child
func (renderer *Renderer) RenderChildCopies(child Child) {
	BindChild(child)
	copies := *(child.GetCopies())
	for x := 0; x < child.GetNumCopies(); x++ {
		renderer.RenderCopy(child, copies[x])
	}
}

// RenderCopy renders a single copy of a child
func (renderer *Renderer) RenderCopy(child Child, c ChildCopy) {
	if renderer.Config.Dimensions == 2 {
		if (child.GetSpecificRenderDistance() != 0 && InBounds2D(c.X, c.Y, float32(renderer.camX), float32(renderer.camY), child.GetSpecificRenderDistance())) ||
			InBounds2D(c.X, c.Y, float32(renderer.camX), float32(renderer.camY), renderer.RenderDistance) {
			child.RenderCopy(c, renderer.MainCamera)
			renderer.DrawChild(child)
			child.AddCurrentCopy(c)
		}
	}
	if renderer.Config.Dimensions == 3 {
		if InBounds3D(c.X, c.Y, c.Z, float32(renderer.camX), float32(renderer.camY), float32(renderer.camZ), renderer.RenderDistance) {
			child.RenderCopy(c, renderer.MainCamera)
			renderer.DrawChild(child)
			child.AddCurrentCopy(c)
		}
	}
}

// BindChild binds the VAO of a child
func BindChild(child Child) {
	gl.BindVertexArray(child.GetVertexArray().id)
	gl.UseProgram(child.GetShaderProgram())
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)

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
	r := Renderer{
		Window:             initGLFW(config),
		ShaderProgram:      initOpenGL(config),
		Children:           []Child{},
		AutomaticRendering: true,
		RenderFunc:         func(r *Renderer) {},
		RenderDistance:     1000,
		Done:               make(chan bool),
		MainCamera:         camera,
		Config:             config,
	}
	r.Window.SetCursorPosCallback(input.MouseCallback)
	r.Window.SetMouseButtonCallback(input.MouseButtonCallback)

	return r
}

func (renderer *Renderer) Initialize(engine *Engine) {
	engine.TextureControl.NewTexture("../rapidengine/border.png", "default")
	dm := NewMaterial(engine.ShaderControl.GetShader("colorLighting"), &engine.Config)
	//dm.BecomeTexture(engine.TextureControl.GetTexture("default"))
	dm.BecomeColor([]float32{0.2, 0.7, 0.4})
	renderer.DefaultMaterial = dm
}

// Instance takes a child and adds it to the renderer's list,
// so that it will be rendered every frame
func (renderer *Renderer) Instance(child Child) {
	renderer.Children = append(renderer.Children, child)
}

// AttachCallback attaches a callback function to the renderer,
// to be called per-frame
func (renderer *Renderer) AttachCallback(f func(*Renderer)) {
	renderer.RenderFunc = f
}

func initGLFW(config *configuration.EngineConfig) *glfw.Window {
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
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}

	window.MakeContextCurrent()

	return window
}

func initOpenGL(config *configuration.EngineConfig) uint32 {
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

	return 0
}

// SetRenderDistance sets the render distance
func (renderer *Renderer) SetRenderDistance(distance float32) {
	renderer.RenderDistance = distance
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
