package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

var MouseX float64
var MouseY float64
var LastMouseX float64
var LastMouseY float64
var LeftMouseButton bool
var RightMouseButton bool

type Input struct {
	Keys             map[string]bool
	MouseX           float64
	MouseY           float64
	LastMouseX       float64
	LastMouseY       float64
	LeftMouseButton  bool
	RightMouseButton bool
}

func MouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	MouseX = xpos
	MouseY = ypos
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	LeftMouseButton = (button == 0 && action == glfw.Press)
	RightMouseButton = (button == 1 && action == glfw.Press)
}

func SwapMousePositions() {
	LastMouseX = MouseX
	LastMouseY = MouseY
}

var KeyMap map[string]glfw.Key = map[string]glfw.Key{
	"w": glfw.KeyW,
	"a": glfw.KeyA,
	"s": glfw.KeyS,
	"d": glfw.KeyD,

	"space":  glfw.KeySpace,
	"shift":  glfw.KeyLeftShift,
	"escape": glfw.KeyEscape,

	"up":    glfw.KeyUp,
	"down":  glfw.KeyDown,
	"left":  glfw.KeyLeft,
	"right": glfw.KeyRight,

	"l": glfw.KeyL,
	"k": glfw.KeyK,
	"o": glfw.KeyO,
	"p": glfw.KeyP,

	"q":    glfw.KeyQ,
	"e":    glfw.KeyE,
	"ctrl": glfw.KeyLeftControl,
}
