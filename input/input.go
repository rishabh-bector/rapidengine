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

type InputControl struct {
	keyMap map[string]glfw.Key
}

func NewInputControl() InputControl {
	return InputControl{keyMap}
}

func (inputControl *InputControl) Update(window *glfw.Window) *Input {
	defer swapMousePositions()
	glfw.PollEvents()
	current := map[string]bool{}
	for name, key := range inputControl.keyMap {
		current[name] = (window.GetKey(key) == glfw.Press)
	}
	return &Input{current, MouseX, MouseY, LastMouseX, LastMouseY, LeftMouseButton, RightMouseButton}
}

func MouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	MouseX = xpos
	MouseY = ypos
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	LeftMouseButton = (button == 0 && action == glfw.Press)
	RightMouseButton = (button == 1 && action == glfw.Press)
}

func swapMousePositions() {
	LastMouseX = MouseX
	LastMouseY = MouseY
}

var keyMap map[string]glfw.Key = map[string]glfw.Key{
	"w":     glfw.KeyW,
	"a":     glfw.KeyA,
	"s":     glfw.KeyS,
	"d":     glfw.KeyD,
	"space": glfw.KeySpace,
	"shift": glfw.KeyLeftShift,
}
