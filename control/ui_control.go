package control

import (
	"rapidengine/input"
	"rapidengine/ui"
)

type UIControl struct {
	Buttons map[string]*ui.Button
}

func NewUIControl() UIControl {
	return UIControl{
		Buttons: make(map[string]*ui.Button),
	}
}

func (ui *UIControl) Update(inputs *input.Input) {
	for _, button := range ui.Buttons {
		button.Update(inputs)
	}
}

func (ui *UIControl) AddButton(button *ui.Button, id string) {
	ui.Buttons[id] = button
}
