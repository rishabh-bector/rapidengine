package terrain

import (
	"rapidengine/child"
	"rapidengine/material"
)

type Foliage struct {
	x float32
	y float32

	width  int
	height int

	spaceX float32
	spaceY float32

	FChild *child.Child3D
}

func NewFoliage(width, height int) Foliage {
	return Foliage{
		width:  width,
		height: height,
		spaceX: 0.1,
		spaceY: 0.1,
	}
}

func (f *Foliage) AttachMaterial(mat *material.FoliageMaterial) {
	f.FChild.AttachMaterial(mat)
	f.FChild.Model.Materials[0] = mat
}
