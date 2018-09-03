package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type InputControl struct {
	keyMap map[string]glfw.Key
}

func NewInputControl() InputControl {
	return InputControl{keyMap}
}

func (inputControl *InputControl) Update(window *glfw.Window) map[string]bool {
	glfw.PollEvents()
	current := map[string]bool{}
	for name, key := range inputControl.keyMap {
		current[name] = (window.GetKey(key) == glfw.Press)
	}
	return current
}

var keyMap map[string]glfw.Key = map[string]glfw.Key{
	"w": glfw.KeyW,
	"a": glfw.KeyA,
	"s": glfw.KeyS,
	"d": glfw.KeyD,
}
