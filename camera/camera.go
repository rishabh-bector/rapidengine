package camera

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Camera interface {
	Look()

	ProcessInput(*glfw.Window)

	GetFirstViewIndex() *float32

	SetPosition(int, int)
	GetPosition() (float32, float32)
}
