package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera2 struct {
	Speed float32

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	View mgl32.Mat4
}

func NewCamera2(position mgl32.Vec3, speed float32) Camera2 {
	return Camera2{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
	}
}

func (camera2 *Camera2) Look() {
	camera2.View = mgl32.LookAtV(
		camera2.Position,
		camera2.Position.Add(camera2.FrontAxis),
		camera2.UpAxis,
	)
}

func (camera2 *Camera2) ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera2.Position = camera2.Position.Add(camera2.UpAxis.Mul(camera2.Speed))
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera2.Position = camera2.Position.Sub(camera2.UpAxis.Mul(camera2.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera2.Position = camera2.Position.Sub(camera2.FrontAxis.Cross(camera2.UpAxis).Normalize().Mul(camera2.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera2.Position = camera2.Position.Add(camera2.FrontAxis.Cross(camera2.UpAxis).Normalize().Mul(camera2.Speed))
	}
}

func (camera2 *Camera2) GetFirstViewIndex() *float32 {
	return &camera2.View[0]
}

func (camera2 *Camera2) ProcessMouse() {
	return
}
