package terrain

import (
	"rapidengine/child"
	"rapidengine/material"
)

type Water struct {
	width  int
	height int

	WChild *child.Child3D
}

func NewWater(width, height int) Water {
	return Water{
		width:  width,
		height: height,
	}
}

func (water *Water) AttachMaterial(mat *material.WaterMaterial) {
	water.WChild.AttachMaterial(mat)
}
