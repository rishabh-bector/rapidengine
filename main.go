package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	log "github.com/sirupsen/logrus"
)

var (
	ScreenWidth  = 1920
	ScreenHeight = 1080

	WindowTitle  = "Test"
	PolygonLines = false

	CameraSpeed       = 0.02
	CameraSensitivity = 0.2

	Dimensions = 2
)

func init() {
	runtime.LockOSThread()
}

func main() {

	camera := NewCamera2D(mgl32.Vec3{0, 0, 0}, float32(CameraSpeed))
	renderer := NewRenderer(render, &camera)
	gl.UseProgram(renderer.ShaderProgram)

	shaders := NewShaders()
	err := shaders.CompileShaders()
	if err != nil {
		log.Fatal(err)
	}

	///   CHILD 1    ///

	child1 := NewChild2D()
	child1.AttachPrimitive(NewRectangle(0.5, 0.5, shaders))
	child1.AttachShader(renderer.ShaderProgram)
	child1.AttachTexture(
		"./texture.png",
		[]float32{
			-1, 1,
			1, 1,
			1, -1,
			-1, -1,
		},
	)
	child1.SetPosition(1, 1)

	renderer.Instance(&child1)

	//////////////////////

	renderer.StartRenderer()
	<-renderer.Done
	shaders.CleanUp()
}

func render(renderer *Renderer) {
	renderer.RenderChildren()
}
