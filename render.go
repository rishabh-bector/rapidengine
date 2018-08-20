package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	log "github.com/sirupsen/logrus"
)

type Renderer struct {
	Window *glfw.Window

	ShaderProgram uint32

	Children []Child

	RenderFunc func(renderer *Renderer)

	MainCamera Camera

	Done chan bool
}

func (renderer *Renderer) startRenderer() {
	for !renderer.Window.ShouldClose() {
		gl.ClearColor(0.5, 0.5, 0.5, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		renderer.RenderFunc(renderer)

		glfw.PollEvents()
		renderer.MainCamera.ProcessInput(renderer.Window)
		renderer.MainCamera.Look()
		renderer.Window.SwapBuffers()
	}
	glfw.Terminate()
	renderer.Done <- true
}

func (renderer *Renderer) RenderChildren() {
	for _, child := range renderer.Children {
		child.PreRender(renderer.MainCamera)
		gl.UseProgram(child.shaderProgram)
		gl.BindVertexArray(child.vertexArray.id)
		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)
		gl.DrawElements(gl.TRIANGLES, child.numVertices, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)
	}
}

func NewRenderer(renderFunc func(renderer *Renderer), camera Camera) Renderer {
	r := Renderer{
		Window:        initGLFW(),
		ShaderProgram: initOpenGL(),
		Children:      []Child{},
		RenderFunc:    renderFunc,
		Done:          make(chan bool),
		MainCamera:    camera,
	}
	r.Window.SetCursorPosCallback(r.MainCamera.ProcessMouse)
	return r
}

func (renderer *Renderer) Instance(child Child) {
	child.PreRender(renderer.MainCamera)
	renderer.Children = append(renderer.Children, child)
}

func initGLFW() *glfw.Window {
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

	window, err := glfw.CreateWindow(ScreenWidth, ScreenHeight, WindowTitle, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
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

	if PolygonLines {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)

	return prog
}

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
