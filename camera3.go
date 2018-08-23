package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var MouseX float64 = float64(ScreenWidth / 2)
var MouseY float64 = float64(ScreenHeight / 2)
var FirstMouse = true

type Camera3 struct {
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

func NewCamera3(position mgl32.Vec3, speed float32) Camera3 {
	return Camera3{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
		Yaw:       0,
		Pitch:     0,
	}
}

func (camera3 *Camera3) Look() {
	camera3.View = mgl32.LookAtV(
		camera3.Position,
		camera3.Position.Add(camera3.FrontAxis),
		camera3.UpAxis,
	)
}

func (camera3 *Camera3) ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera3.Position = camera3.Position.Add(camera3.FrontAxis.Mul(camera3.Speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera3.Position = camera3.Position.Sub(camera3.FrontAxis.Mul(camera3.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera3.Position = camera3.Position.Sub(camera3.FrontAxis.Cross(camera3.UpAxis).Normalize().Mul(camera3.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera3.Position = camera3.Position.Add(camera3.FrontAxis.Cross(camera3.UpAxis).Normalize().Mul(camera3.Speed))
	}
}

func (camera3 *Camera3) GetFirstViewIndex() *float32 {
	return &camera3.View[0]
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

func (camera3 *Camera3) ProcessMouse() {
	if FirstMouse {
		camera3.MouseLastX = MouseX
		camera3.MouseLastY = MouseY
		FirstMouse = false
	}
	xOffset := MouseX - camera3.MouseLastX
	yOffset := camera3.MouseLastY - MouseY
	camera3.MouseLastX = MouseX
	camera3.MouseLastY = MouseY
	xOffset *= CameraSensitivity
	yOffset *= CameraSensitivity
	camera3.Yaw += float32(xOffset)
	camera3.Pitch += float32(yOffset)
	if camera3.Pitch > 89 {
		camera3.Pitch = 89
	}
	if camera3.Pitch < -89 {
		camera3.Pitch = -89
	}
	camera3.FrontAxis = CalculateDirection(camera3.Pitch, camera3.Yaw).Normalize()
}
