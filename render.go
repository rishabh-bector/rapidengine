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
)

// Renderer contains the information required for
// the main render loop
type Renderer struct {
	Window *glfw.Window

	ShaderProgram uint32

	Children []Child

	RenderFunc func(renderer *Renderer)

	MainCamera Camera

	Config *EngineConfig

	Done chan bool
}

// StartRenderer starts the main render loop
func (renderer *Renderer) StartRenderer() {
	renderer.PreRenderChildren()

	for !renderer.Window.ShouldClose() {
		gl.ClearColor(0.5, 0.5, 0.5, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		renderer.RenderFunc(renderer)

		glfw.PollEvents()
		renderer.MainCamera.ProcessInput(renderer.Window)
		renderer.MainCamera.Look()
		renderer.Window.SwapBuffers()
		CheckError("loop")
	}
	glfw.Terminate()
	renderer.Done <- true
}

// PreRenderChildren calls the PreRender method of each child,
// for initialization
func (renderer *Renderer) PreRenderChildren() {
	for _, child := range renderer.Children {
		child.PreRender(renderer.MainCamera)
	}
}

// RenderChildren binds the appropriate shaders and Vertex Array for each child,
// and draws them to the screen using an element buffer
func (renderer *Renderer) RenderChildren() {
	for _, child := range renderer.Children {
		// Call the child's update method to update transform matrices
		child.Update(renderer.MainCamera)

		// Bind Shader Program & Vertex Array
		gl.UseProgram(child.GetShaderProgram())
		gl.BindVertexArray(child.GetVertexArray().id)
		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)

		// Bind child's texture
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, child.GetTexture())

		// Draw elements and unbind array
		gl.DrawElements(gl.TRIANGLES, child.GetNumVertices(), gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)
	}
}

// NewRenderer creates a new renderer, and takes in a renderFunc which
// is called every frame, allowing the User to have frame-by-frame control
func NewRenderer(camera Camera, config *EngineConfig) Renderer {
	r := Renderer{
		Window:        initGLFW(config),
		ShaderProgram: initOpenGL(config),
		Children:      []Child{},
		RenderFunc:    func(r *Renderer) {},
		Done:          make(chan bool),
		MainCamera:    camera,
		Config:        config,
	}
	//r.Window.SetCursorPosCallback(MouseCallback)
	return r
}

// Instance takes a child and adds it to the renderer's list,
// so that it will be rendered every frame
func (renderer *Renderer) Instance(child Child) {
	child.PreRender(renderer.MainCamera)
	renderer.Children = append(renderer.Children, child)
}

// AttachCallback attaches a callback function to the renderer,
// to be called per-frame
func (renderer *Renderer) AttachCallback(f func(*Renderer)) {
	renderer.RenderFunc = f
}

func initGLFW(config *EngineConfig) *glfw.Window {
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

	window, err := glfw.CreateWindow(config.ScreenWidth, config.ScreenHeight, config.WindowTitle, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	//window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.MakeContextCurrent()

	return window
}

func initOpenGL(config *EngineConfig) uint32 {
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("Using OpenGL Version ", version)

	vertexShader, err := CompileShader(VertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatal(err)
	}

	fragmentShader, err := CompileShader(FragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatal(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	if config.PolygonLines {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)

	return prog
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
