package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Speed float32

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	Pitch float32
	Yaw   float32

	MouseLastX float64
	MouseLastY float64

	View mgl32.Mat4
}

func NewCamera(position mgl32.Vec3, speed float32) Camera {
	return Camera{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
	}
}

func (camera *Camera) Look() {
	camera.View = mgl32.LookAtV(
		camera.Position,
		camera.Position.Add(camera.FrontAxis),
		camera.UpAxis,
	)
}

func (camera *Camera) ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.Position = camera.Position.Add(camera.FrontAxis.Mul(camera.Speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.Position = camera.Position.Sub(camera.FrontAxis.Mul(camera.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.Position = camera.Position.Sub(camera.FrontAxis.Cross(camera.UpAxis).Normalize().Mul(camera.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.Position = camera.Position.Add(camera.FrontAxis.Cross(camera.UpAxis).Normalize().Mul(camera.Speed))
	}
}

func CalculateDirection(pitch, yaw float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Cos(float64(mgl32.DegToRad(yaw)))),
		float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Sin(float64(mgl32.DegToRad(yaw)))),
	}
}

func (camera *Camera) ProcessMouse(w *glfw.Window, xpos float64, ypos float64) {
	xOffset := xpos - camera.MouseLastX
	yOffset := camera.MouseLastY - ypos
	camera.MouseLastX = xpos
	camera.MouseLastY = ypos
	xOffset *= CameraSensitivity
	yOffset *= CameraSensitivity
	camera.Yaw += float32(xOffset)
	camera.Pitch += float32(yOffset)
	if camera.Pitch > 89 {
		camera.Pitch = 89
	}
	if camera.Pitch < -89 {
		camera.Pitch = -89
	}
	camera.FrontAxis = CalculateDirection(camera.Pitch, camera.Yaw).Normalize()
}
