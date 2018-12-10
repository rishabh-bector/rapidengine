package camera

import (
	"rapidengine/input"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	Look()

	DefaultControls(*input.Input)
	MoveUp()
	MoveDown()
	MoveLeft()
	MoveRight()
	MoveForward()
	MoveBackward()

	ChangeYaw(float32)
	ChangePitch(float32)
	ChangeRoll(float32)

	GetFirstViewIndex() *float32
	GetStaticView() mgl32.Mat4

	SetPosition(float32, float32, float32)
	GetPosition() (float32, float32, float32)
	SetSpeed(float32)

	ProcessMouse(float64, float64, float64, float64)
}
