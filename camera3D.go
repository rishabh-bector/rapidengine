package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var MouseX float64 = float64(ScreenWidth / 2)
var MouseY float64 = float64(ScreenHeight / 2)
var FirstMouse = true

type Camera3D struct {
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

func NewCamera3D(position mgl32.Vec3, speed float32) Camera3D {
	return Camera3D{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
		Yaw:       0,
		Pitch:     0,
	}
}

func (camera3D *Camera3D) Look() {
	camera3D.View = mgl32.LookAtV(
		camera3D.Position,
		camera3D.Position.Add(camera3D.FrontAxis),
		camera3D.UpAxis,
	)
}

func (camera3D *Camera3D) ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera3D.Position = camera3D.Position.Add(camera3D.FrontAxis.Mul(camera3D.Speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera3D.Position = camera3D.Position.Sub(camera3D.FrontAxis.Mul(camera3D.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera3D.Position = camera3D.Position.Sub(camera3D.FrontAxis.Cross(camera3D.UpAxis).Normalize().Mul(camera3D.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera3D.Position = camera3D.Position.Add(camera3D.FrontAxis.Cross(camera3D.UpAxis).Normalize().Mul(camera3D.Speed))
	}
}

func (camera3D *Camera3D) GetFirstViewIndex() *float32 {
	return &camera3D.View[0]
}

func CalculateDirection(pitch, yaw float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Cos(float64(mgl32.DegToRad(yaw)))),
		float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Sin(float64(mgl32.DegToRad(yaw)))),
	}
}

func MouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	MouseX = xpos
	MouseY = ypos
}

func (camera3D *Camera3D) ProcessMouse() {
	if FirstMouse {
		camera3D.MouseLastX = MouseX
		camera3D.MouseLastY = MouseY
		FirstMouse = false
	}
	xOffset := MouseX - camera3D.MouseLastX
	yOffset := camera3D.MouseLastY - MouseY
	camera3D.MouseLastX = MouseX
	camera3D.MouseLastY = MouseY
	xOffset *= CameraSensitivity
	yOffset *= CameraSensitivity
	camera3D.Yaw += float32(xOffset)
	camera3D.Pitch += float32(yOffset)
	if camera3D.Pitch > 89 {
		camera3D.Pitch = 89
	}
	if camera3D.Pitch < -89 {
		camera3D.Pitch = -89
	}
	camera3D.FrontAxis = CalculateDirection(camera3D.Pitch, camera3D.Yaw).Normalize()
}
