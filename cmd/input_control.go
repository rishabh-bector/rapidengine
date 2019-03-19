package cmd

import (
	"rapidengine/input"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type InputControl struct {
	keyMap map[string]glfw.Key
}

func NewInputControl() InputControl {
	return InputControl{input.KeyMap}
}

func (inputControl *InputControl) Update(window *glfw.Window) *input.Input {
	defer input.SwapMousePositions()
	glfw.PollEvents()
	current := map[string]bool{}
	for name, key := range inputControl.keyMap {
		current[name] = (window.GetKey(key) == glfw.Press)
	}
	return &input.Input{
		current,

		input.MouseX,
		input.MouseY,

		input.LastMouseX,
		input.LastMouseY,

		input.LeftMouseButton,
		input.RightMouseButton,
		input.MiddleMouseButton,

		input.ScrollXOff,
		input.ScrollYOff,
		input.Scroll,
	}
}
