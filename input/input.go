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
var MiddleMouseButton bool

var ScrollXOff float64
var ScrollYOff float64
var Scroll float64

type Input struct {
	Keys map[string]bool

	MouseX float64
	MouseY float64

	LastMouseX float64
	LastMouseY float64

	LeftMouseButton   bool
	RightMouseButton  bool
	MiddleMouseButton bool

	ScrollX float64
	ScrollY float64
	Scroll  float64
}

func MouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	MouseX = xpos
	MouseY = ypos
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	LeftMouseButton = (button == 0 && action == glfw.Press)
	RightMouseButton = (button == 1 && action == glfw.Press)
	MiddleMouseButton = (button == 2 && action == glfw.Press)
}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	ScrollXOff = xoff
	ScrollYOff = yoff
	Scroll += yoff
}

func SwapMousePositions() {
	LastMouseX = MouseX
	LastMouseY = MouseY
}

var KeyMap map[string]glfw.Key = map[string]glfw.Key{
	"space":  glfw.KeySpace,
	"shift":  glfw.KeyLeftShift,
	"escape": glfw.KeyEscape,

	"up":    glfw.KeyUp,
	"down":  glfw.KeyDown,
	"left":  glfw.KeyLeft,
	"right": glfw.KeyRight,

	"a": glfw.KeyA,
	"b": glfw.KeyB,
	"c": glfw.KeyC,
	"d": glfw.KeyD,
	"e": glfw.KeyE,
	"f": glfw.KeyF,
	"g": glfw.KeyG,
	"h": glfw.KeyH,
	"i": glfw.KeyI,
	"j": glfw.KeyJ,
	"k": glfw.KeyK,
	"l": glfw.KeyL,
	"m": glfw.KeyM,
	"n": glfw.KeyN,
	"o": glfw.KeyO,
	"p": glfw.KeyP,
	"q": glfw.KeyQ,
	"r": glfw.KeyR,
	"s": glfw.KeyS,
	"t": glfw.KeyT,
	"u": glfw.KeyU,
	"v": glfw.KeyV,
	"w": glfw.KeyW,
	"x": glfw.KeyX,
	"y": glfw.KeyY,
	"z": glfw.KeyZ,

	"ctrl_left":  glfw.KeyLeftControl,
	"ctrl_right": glfw.KeyRightControl,
}
