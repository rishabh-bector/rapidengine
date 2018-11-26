package ui

import (
	"rapidengine/child"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/input"
)

type ProgressBar struct {
	BackChild *child.Child2D
	BarChild  *child.Child2D

	TextBxLeft  *TextBox
	TextBxRight *TextBox

	transform geometry.Transform

	barScaleX float32
	barScaleY float32

	currentPercentage float32
}

func NewProgressBar(config *configuration.EngineConfig) ProgressBar {
	pb := ProgressBar{
		currentPercentage: 50,
		transform:         geometry.NewTransform(0, 0, 0, 100, 25, 0),
		barScaleX:         1,
		barScaleY:         1,
	}

	return pb
}

func (pb *ProgressBar) Initialize() {
	pb.SetPosition(pb.transform.X, pb.transform.Y)
	pb.SetDimensions(pb.transform.SX, pb.transform.SY)

	pb.BackChild.Static = true
	pb.BarChild.Static = true
}

func (pb *ProgressBar) SetBarScale(sx, sy float32) {
	pb.barScaleX = sx
	pb.barScaleY = sy
	pb.Initialize()
}

func (pb *ProgressBar) AttachText(left *TextBox, right *TextBox) {
	pb.TextBxLeft = left
	pb.TextBxRight = right
}

func (pb *ProgressBar) IncrementPercentage(delta float32) {
	pb.SetPercentage(pb.currentPercentage + delta)
}

func (pb *ProgressBar) SetPercentage(percent float32) {
	if percent > 100 || percent < 0 {
		return
	}
	pb.currentPercentage = percent
}

func (pb *ProgressBar) GetPercentage() float32 {
	return pb.currentPercentage
}

//  --------------------------------------------------
//  Interface
//  --------------------------------------------------

func (pb *ProgressBar) Update(inputs *input.Input) {
	pb.BarChild.ScaleX = (pb.currentPercentage / 100) * (pb.transform.SX * pb.barScaleX)
}

func (pb *ProgressBar) SetPosition(x, y float32) {
	pb.transform.X = x
	pb.transform.Y = y

	pb.BackChild.X = x
	pb.BackChild.Y = y

	pb.BarChild.X = x + (pb.transform.SX-(pb.transform.SX*pb.barScaleX))/2
	pb.BarChild.Y = y + (pb.transform.SY-(pb.transform.SY*pb.barScaleY))/2
}

func (pb *ProgressBar) SetDimensions(width, height float32) {
	pb.transform.SX = width
	pb.transform.SY = height

	pb.BackChild.ScaleX = pb.transform.SX
	pb.BackChild.ScaleY = pb.transform.SY

	pb.BarChild.ScaleX = pb.transform.SX * pb.barScaleX
	pb.BarChild.ScaleY = pb.transform.SY * pb.barScaleY
}

func (pb *ProgressBar) GetTransform() geometry.Transform {
	return pb.transform
}

func (pb *ProgressBar) GetChildren() []*child.Child2D {
	return []*child.Child2D{pb.BackChild, pb.BarChild}
}

func (pb *ProgressBar) GetTextBoxes() []*TextBox {
	return []*TextBox{pb.TextBxLeft, pb.TextBxRight}
}
