package rapidengine

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera2D struct {
	Speed float32

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	View mgl32.Mat4
}

func NewCamera2D(position mgl32.Vec3, speed float32) *Camera2D {
	return &Camera2D{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
	}
}

func (camera2D *Camera2D) Look() {
	camera2D.View = mgl32.LookAtV(
		camera2D.Position,
		camera2D.Position.Add(camera2D.FrontAxis),
		camera2D.UpAxis,
	)
}

func (camera2D *Camera2D) ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera2D.Position = camera2D.Position.Add(camera2D.UpAxis.Mul(camera2D.Speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera2D.Position = camera2D.Position.Sub(camera2D.UpAxis.Mul(camera2D.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera2D.Position = camera2D.Position.Sub(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera2D.Position = camera2D.Position.Add(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
	}
}

func (camera2D *Camera2D) GetFirstViewIndex() *float32 {
	return &camera2D.View[0]
}

func (camera2D *Camera2D) ProcessMouse() {
	return
}
