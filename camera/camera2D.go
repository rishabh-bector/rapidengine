package camera

import (
	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/configuration"
)

type Camera2D struct {
	Speed float32

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	View mgl32.Mat4

	config *configuration.EngineConfig
}

func NewCamera2D(position mgl32.Vec3, speed float32, config *configuration.EngineConfig) *Camera2D {
	return &Camera2D{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
		config:    config,
	}
}

func (camera2D *Camera2D) Look() {
	camera2D.View = mgl32.LookAtV(
		camera2D.Position,
		camera2D.Position.Add(camera2D.FrontAxis),
		camera2D.UpAxis,
	)
}

func (camera2D *Camera2D) MoveUp() {
	camera2D.Position = camera2D.Position.Add(camera2D.UpAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveDown() {
	camera2D.Position = camera2D.Position.Sub(camera2D.UpAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveLeft() {
	camera2D.Position = camera2D.Position.Sub(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveRight() {
	camera2D.Position = camera2D.Position.Add(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveForward() {
	camera2D.Position = camera2D.Position.Add(camera2D.FrontAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveBackward() {
	camera2D.Position = camera2D.Position.Sub(camera2D.FrontAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) ChangeYaw(y float32)   {}
func (camera2D *Camera2D) ChangePitch(p float32) {}

func (camera2D *Camera2D) GetFirstViewIndex() *float32 {
	return &camera2D.View[0]
}

func (camera2D *Camera2D) GetPosition() (float32, float32, float32) {
	return ((camera2D.Position.X() / 2) * float32(camera2D.config.ScreenWidth)) + float32(camera2D.config.ScreenWidth/2),
		((camera2D.Position.Y() / 2) * float32(camera2D.config.ScreenHeight)) + float32(camera2D.config.ScreenHeight/2), 0
}

func (camera2D *Camera2D) SetPosition(x, y, z float32) {
	camera2D.Position = mgl32.Vec3{
		(x - float32(camera2D.config.ScreenWidth/2)) / float32(camera2D.config.ScreenWidth/2),
		(y - float32(camera2D.config.ScreenHeight/2)) / float32(camera2D.config.ScreenHeight/2),
		camera2D.Position.Z(),
	}
}

func (camera2D *Camera2D) SetSpeed(s float32) {
	camera2D.Speed = s
}

func (camera2D *Camera2D) ProcessMouse(mx, my, lmx, lmy float64) {
	return
}
