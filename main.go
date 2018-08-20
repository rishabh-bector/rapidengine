package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	log "github.com/sirupsen/logrus"
)

var (
	ScreenWidth  = 1440
	ScreenHeight = 900

	WindowTitle  = "Test"
	PolygonLines = false

	CameraSpeed       = 0.02
	CameraSensitivity = 0.05
)

func init() {
	runtime.LockOSThread()
}

func main() {

	renderer := NewRenderer(render, NewCamera(mgl32.Vec3{0, 0, 3}, float32(CameraSpeed)))
	gl.UseProgram(renderer.ShaderProgram)

	shaders := NewShaders()
	err := shaders.CompileShaders()
	if err != nil {
		log.Fatal(err)
	}

	///   CHILD 1    ///

	child1 := NewChild()
	child1.AttachPrimitive(NewCube(shaders))
	child1.AttachShader(renderer.ShaderProgram)
	child1.AttachTexture(
		"./texture.png",
		CubeTextures,
	)

	renderer.Instance(child1)

	//////////////////////

	renderer.startRenderer()
	<-renderer.Done
	shaders.CleanUp()
}

func render(renderer *Renderer) {
	renderer.RenderChildren()
}
