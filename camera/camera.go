package camera

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Camera interface {
	Look()

	ProcessInput(*glfw.Window)

	GetFirstViewIndex() *float32

	SetPosition(float32, float32)
	GetPosition() (float32, float32)
	SetSpeed(float32)
}
