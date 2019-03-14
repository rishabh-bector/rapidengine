package terrain

import (
	"rapidengine/child"
	"rapidengine/material"
)

type Terrain struct {
	width  int
	height int

	TChild *child.Child3D
}

func NewTerrain(width, height int) Terrain {
	return Terrain{
		width:  width,
		height: height,
	}
}

func (terrain *Terrain) AttachMaterial(mat *material.TerrainMaterial) {
	terrain.TChild.Model.Materials[0] = mat
}
