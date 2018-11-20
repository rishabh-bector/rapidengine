package ui

import (
	"rapidengine/geometry"
	"rapidengine/input"
)

type Element interface {
	Update(inputs *input.Input)

	SetPosition(x, y float32)
	SetDimensions(width, height float32)

	GetTransform() geometry.Transform
}
