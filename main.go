package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	log "github.com/sirupsen/logrus"
)

var (
	ScreenWidth  = 1440
	ScreenHeight = 900

	WindowTitle  = "Test"
	PolygonLines = false
)

func init() {
	runtime.LockOSThread()
}

func main() {
	window := initGLFW()
	defer glfw.Terminate()

	program := initOpenGL()
	gl.UseProgram(program)
	CheckError("initOpenGL")

	renderer := Renderer{
		window:        window,
		shaderProgram: program,
		children:      []Child{},
	}

	shaders := NewShaders()
	err := shaders.CompileShaders()
	if err != nil {
		log.Fatal(err)
	}

	UnbindBuffers()

	///   CHILD 1    ///

	child1 := NewChild()
	child1.AttachPrimitive(NewCube(shaders))
	child1.AttachShader(program)
	child1.AttachTexture(
		"./texture.png",
		CubeTextures,
	)

	renderer.Instance(child1)

	//////////////////////

	for !window.ShouldClose() {
		gl.ClearColor(0.5, 0.5, 0.5, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		renderer.RenderChildren()

		glfw.PollEvents()
		renderer.window.SwapBuffers()
	}

	shaders.CleanUp()
}
